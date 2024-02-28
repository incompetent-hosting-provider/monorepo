package terraformbinary

import (
	"context"
	"os"
	"path/filepath"

	version "github.com/hashicorp/go-version"
	product "github.com/hashicorp/hc-install/product"
	releases "github.com/hashicorp/hc-install/releases"
	tfexec "github.com/hashicorp/terraform-exec/tfexec"
	log "github.com/rs/zerolog/log"
)

func installTerraform(binary_dir string, terraform_version string) (string, error) {
	log.Debug().Msgf("Trying to install terraform binary version %s", terraform_version)
	installer := &releases.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion(terraform_version)),
		InstallDir: binary_dir,
	}

	execPath, err := installer.Install(context.TODO())
	if err != nil {
		return "", err
	}

	log.Debug().Msgf("Installed terraform binary to %s", execPath)

	return execPath, nil
}

func ensureTerraform(binary_dir string, terraform_version string, terraform_dir string) (*tfexec.Terraform, error) {
	log.Debug().Msgf("Ensuring terraform version %s", terraform_version)
	execPath := filepath.Join(binary_dir, "terraform")

	// Check if the binary directory exists
	if _, err := os.Stat(binary_dir); os.IsNotExist(err) {
		log.Debug().Msgf("Binary directory %s not found", binary_dir)
		log.Info().Msg("Binary directory not found, creating...")
		err := os.Mkdir(binary_dir, 0755)
		if err != nil {
			return nil, err
		}
	}

	// Check if the terraform binary is already present
	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		log.Debug().Msgf("Terraform binary not found at %s", execPath)
		log.Info().Msg("Terraform binary not found, downloading...")
		execPath, err = installTerraform(binary_dir, terraform_version)
		if err != nil {
			return nil, err
		}

	} else {
		log.Debug().Msgf("Terraform binary found at %s", execPath)
		log.Info().Msg("Terraform binary found, skipping download...")

		log.Debug().Msg("Creating new terraform struct")
		tf, err := tfexec.NewTerraform(terraform_dir, execPath)
		if err != nil {
			return nil, err
		}

		log.Debug().Msg("Trying to get terraform version")
		tfVersion, _, err := tf.Version(context.TODO(), true)
		if err != nil {
			return nil, err
		}

		if tfVersion.String() != terraform_version {
			log.Info().Msgf("Terraform version mismatch, %s != %s downloading version %s", tfVersion.String(), terraform_version, terraform_version)
			execPath, err = installTerraform(binary_dir, terraform_version)
			if err != nil {
				return nil, err
			}

		}

		log.Debug().Msgf("Trying to get new terraform version")
		tfVersion, _, err = tf.Version(context.TODO(), true)
		if err != nil {
			return nil, err
		}

		log.Debug().Msgf("Terraform installed in version: %s", tfVersion)
	}

	tf, err := tfexec.NewTerraform(terraform_dir, execPath)
	if err != nil {
		return nil, err
	}

	return tf, nil
}
