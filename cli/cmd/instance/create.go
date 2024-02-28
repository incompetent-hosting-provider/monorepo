package instance

import (
	"cli/internal/authentication"
	"cli/internal/backend"
	"cli/internal/messages"
	"cli/internal/models"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// Instances Create Command
//
// Runs the create instance prompt so the user can create a new instance
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a predefined or custom instance.",
	Long:  "Creates a predefined or custom instance.",
	Run: func(cmd *cobra.Command, args []string) {
		tokens := authentication.GetCurrentAuthentication()
		if tokens == nil {
			messages.DisplayNotLoggedInMessage()
			return
		}

		isCustomInstance, err := promptIsCustomInstance()
		if err != nil {
			handleCreateError(err)
			return
		}

		var preset models.InstancePreset
		var image models.ContainerImage
		if isCustomInstance {
			image, err = promptImageCustom()
			if err != nil {
				handleCreateError(err)
				return
			}
		} else {
			preset, err = promptImageFromTemplates(tokens.AccessToken)
			if err != nil {
				handleCreateError(err)
				return
			}
		}

		name, err := promptInstanceName()
		if err != nil {
			handleCreateError(err)
			return
		}

		description, err := promptInstanceDescription()
		if err != nil {
			handleCreateError(err)
			return
		}

		if isCustomInstance {
			resp, err := createCustomInstance(tokens.AccessToken, name, description, image)
			if err != nil {
				handleCreateError(err)
				return
			}

			fmt.Printf("Instance created with ID: %s\n", resp.InstanceID)
			fmt.Println("The instance will be available shortly.")
		} else {
			resp, err := createTemplateInstance(tokens.AccessToken, name, description, preset.ID)
			if err != nil {
				handleCreateError(err)
				return
			}

			fmt.Printf("Instance created with ID: %s\n", resp.InstanceID)
			if resp.EnvVars != nil && len(resp.EnvVars) > 0 {
				fmt.Println("The following environment variables were set:")
				for key, value := range resp.EnvVars {
					fmt.Printf("%s: %s\n", key, value)
				}
			}
			fmt.Println("The instance will be available shortly.")
		}
	},
}

func createCustomInstance(token authentication.AccessToken, name string, description string, image models.ContainerImage) (backend.CreateCustomInstanceResponse, error) {
	openPorts, err := promptInstanceOpenPorts()
	if err != nil {
		return backend.CreateCustomInstanceResponse{}, err
	}

	envVariables, err := promptEnvVariables()
	if err != nil {
		return backend.CreateCustomInstanceResponse{}, err
	}

	request := backend.CreateCustomInstanceRequest{
		Name: name,
		Description: description,
		Image: image,
		Ports: openPorts,
		EnvVars: envVariables,
	}

	return backend.DefaultBackendClient.CreateCustomInstance(token, request, true)
}

func createTemplateInstance(token authentication.AccessToken, name string, description string, preset int) (backend.CreatePresetInstanceResponse, error){
	request := backend.CreatePresetInstanceRequest{
		PresetID: preset,
		Name: name,
		Description: description,
	}

	return backend.DefaultBackendClient.CreatePresetInstance(token, request, true)
}

func handleCreateError(err error) {
	if errors.Is(err, backend.ErrNotAuthenticated) {
		messages.DisplaySessionExpiredMessage()
	} else {
		fmt.Printf("Creation of instance failed: %s\n", err.Error())
		fmt.Println("Please try again later.")
	}
}

func promptIsCustomInstance() (bool, error) {
	prompt := promptui.Select{
		Label: "Select Instance Type",
		Items: []string{"Predefined", "Custom"},
	}

	result, _, err := prompt.Run()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func promptImageFromTemplates(token authentication.AccessToken) (models.InstancePreset, error){
	templates, err := backend.DefaultBackendClient.GetInstancePresets(token, true)
	if err != nil {
		return models.InstancePreset{}, err
	}

	prompt := promptui.Select{
		Label: "Select Template",
		Items: templates,
	}

	selected, _, err := prompt.Run()

	return templates[selected], err
}

func promptImageCustom() (models.ContainerImage, error) {
	imageIdPrompt := promptui.Prompt{
		Label: "Image ID (from dockerhub)",
		Validate: func(input string) error {
			if strings.Trim(input, " ") == "" {
				return errors.New("image cannot be empty")
			}
			return nil
		},
	}

	imageId, err := imageIdPrompt.Run()
	if err != nil {
		return models.ContainerImage{}, err
	}

	imageVersionPrompt := promptui.Prompt{
		Label: "Image Version (empty for latest)",
		Validate: func(input string) error {
			return nil
		},
	}

	imageVersion, err := imageVersionPrompt.Run()
	if err != nil {
		return models.ContainerImage{}, err
	}

	if imageVersion == "" {
		imageVersion = "latest"
	}

	return models.ContainerImage{
		Name: imageId,
		Version: imageVersion,
	}, nil
}

// Prompts the user for the instance name
func promptInstanceName() (string, error) {
	prompt := promptui.Prompt{
		Label: "Instance Name",
		Validate: func(input string) error {
			trimmedInput := strings.Trim(input, " ")
			if len(trimmedInput) < 3 {
				return errors.New("instance name must be at least 3 characters long")
			} else if len(trimmedInput) > 50 {
				return errors.New("instance name must be at most 50 characters long")
			}
			return nil
		},
	}

	return prompt.Run()
}

// Prompts the user for the instance description
func promptInstanceDescription() (string, error) {
	prompt := promptui.Prompt{
		Label: "Instance Description",
		Validate: func(input string) error {
			trimmedInput := strings.Trim(input, " ")
			if len(trimmedInput) > 200 {
				return errors.New("instance description must be at most 200 characters long")
			}
			return nil
		},
	}

	return prompt.Run()
}

// Prompts the user for the instance open ports
func promptInstanceOpenPorts() ([]int, error) {
	fmt.Println("Please enter the ports you want to open for your instance (empty to finish)")
	prompt := promptui.Prompt{
		Label: "Open Port",
		Validate: func(input string) error {
			if input == "" {
				return nil
			}

			port, err := strconv.Atoi(input)
			if err != nil {
				return err
			}

			if port < 0 || port > 65535 {
				return errors.New("port must be between 0 and 65535")
			}

			return nil
		},
	}

	ports := make([]int, 0)
	for {
		result, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		if result == "" {
			break
		}

		port, err := strconv.Atoi(result)
		if err != nil {
			return nil, err
		}

		ports = append(ports, port)
	}

	return ports, nil
}

// Prompts the user for the instance environment variables
func promptEnvVariables() (map[string]string, error) {
	fmt.Println("Please define the environment variables for your instance")
	prompt := promptui.Prompt{
		Validate: func(input string) error {
			return nil
		},
	}

	envVariables := make(map[string]string)
	for {
		prompt.Label = "Environment Variable Name (empty to finish)"
		key, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		if key == "" {
			break
		}

		prompt.Label = fmt.Sprintf("Value for <<%s>>", key)
		value, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		envVariables[key] = value
	}

	return envVariables, nil
}