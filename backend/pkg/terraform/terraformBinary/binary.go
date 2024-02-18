package terraformbinary

import (
	"context"
	"incompetent-hosting-provider/backend/pkg/util"
	"os"

	version "github.com/hashicorp/go-version"
	product "github.com/hashicorp/hc-install/product"
	releases "github.com/hashicorp/hc-install/releases"
	tfexec "github.com/hashicorp/terraform-exec/tfexec"
	log "github.com/rs/zerolog/log"
)

func installTerraform(binary_dir string, terraform_version string) (string, error) {
	installer := &releases.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion(terraform_version)),
		InstallDir: binary_dir,
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		return "", err
	}

	log.Debug().Msgf("Installed terraform binary to %s", execPath)

	return execPath, nil
}

func ensureTerraform(binary_dir string, terraform_version string, terraform_dir string) (*tfexec.Terraform, error) {
	execPath := ""

	// Check if the terraform binary is already present
	if _, err := os.Stat(binary_dir + "/terraform"); os.IsNotExist(err) {

		log.Debug().Msg("Terraform binary not found, downloading...")
		execPath, err = installTerraform(binary_dir, terraform_version)
		if err != nil {
			return nil, err
		}

	} else {

		log.Debug().Msg("Terraform binary found, skipping download...")
		execPath = binary_dir + "/terraform"

		util.CreateDirFromAbolutePathIfNotExist(binary_dir, 0700)

		tf, err := tfexec.NewTerraform(terraform_dir, execPath)
		if err != nil {
			return nil, err
		}

		tfVersion, _, err := tf.Version(context.Background(), true)
		if err != nil {
			return nil, err
		}

		if tfVersion.String() != terraform_version {

			log.Debug().Msgf("Terraform version mismatch, %s != %s downloading version %s", tfVersion.String(), terraform_version, terraform_version)
			execPath, err = installTerraform(binary_dir, terraform_version)
			if err != nil {
				return nil, err
			}

		}

		tfVersion, _, err = tf.Version(context.Background(), true)
		if err != nil {
			return nil, err
		}

		log.Debug().Msgf("Terraform version: %s", tfVersion)
	}

	tf, err := tfexec.NewTerraform(terraform_dir, execPath)
	if err != nil {
		return nil, err
	}

	tf.Version(context.Background(), true)

	return tf, nil
}
