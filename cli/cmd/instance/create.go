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
			displayErrorMessage(err)
			return
		}

		var preset models.InstancePreset
		var image models.ContainerImage
		if isCustomInstance {
			image, err = promptImageCustom()
			if err != nil {
				displayErrorMessage(err)
				return
			}
		} else {
			preset, err = promptImageFromTemplates()
			if err != nil {
				displayErrorMessage(err)
				return
			}
		}

		name, err := promptInstanceName()
		if err != nil {
			displayErrorMessage(err)
			return
		}

		description, err := promptInstanceDescription()
		if err != nil {
			displayErrorMessage(err)
			return
		}

		var createdId string
		var createError error
		if isCustomInstance {
			openPorts, err := promptInstanceOpenPorts()
			if err != nil {
				displayErrorMessage(err)
				return
			}
	
			envVariables, err := promptEnvVariables()
			if err != nil {
				displayErrorMessage(err)
				return
			}

			request := backend.CreateCustomInstanceRequest{
				Name: name,
				Description: description,
				Image: image,
				Ports: openPorts,
				EnvVars: envVariables,
			}

			createdId, createError = backend.DefaultBackendClient.CreateCustomInstance(tokens.AccessToken, request, true)
		} else {
			request := backend.CreatePresetInstanceRequest{
				PresetID: preset.ID,
				Name: name,
				Description: description,
			}

			createdId, createError = backend.DefaultBackendClient.CreatePresetInstance(tokens.AccessToken, request, true)
		}

		if errors.Is(createError, backend.ErrNotAuthenticated) {
			messages.DisplaySessionExpiredMessage()
			return
		} else if createError != nil {
			displayErrorMessage(createError)
			return
		}

		fmt.Printf("Instance %s created successfully!\n", name)
		fmt.Printf("Created instance id: %s\n", createdId)
	},
}

func displayErrorMessage(err error) {
	fmt.Printf("Creation of instance failed: %s\n", err.Error())
	fmt.Println("Please try again later.")
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

func promptImageFromTemplates() (models.InstancePreset, error){
	templates, err := backend.DefaultBackendClient.GetInstancePresets()
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