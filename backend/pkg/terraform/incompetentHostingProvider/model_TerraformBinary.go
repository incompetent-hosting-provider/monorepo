package incompetenthostingprovider

import (
	"context"
	tfbin "incompetent-hosting-provider/backend/pkg/terraform/terraformBinary"

	tfexec "github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

type TerraformBinary struct {
	tfbin *tfbin.TerraformBinary
}

func NewTfBin(install_directory string, required_version string, working_directory string, environment_files []string) *TerraformBinary {
	return &TerraformBinary{
		tfbin: tfbin.NewTfBin(install_directory, required_version, working_directory, environment_files),
	}
}

func (i *TerraformBinary) GetInstance() (*tfexec.Terraform, error) {
	return i.tfbin.GetInstance()
}

func (i *TerraformBinary) GetEnvironmentFiles() []string {
	return i.tfbin.GetEnvironmentFiles()
}

func (i *TerraformBinary) ApplyTerraform(containerToCreate []DockerMySQL, containerToDestroy []DockerMySQL) (*tfjson.State, error) {
	tf, err := i.GetInstance()
	if err != nil {
		return nil, err
	}

	ihpenv, err := ReadCustomEnv("./TerraDocker/terraform.tfvars.json")
	if err != nil {
		return nil, err
	}
	// This is a quick fix, since the envirnment files are also used in the tfbin where the CWD is another directory
	ihpcredsenv, err := ReadCredsEnv("./TerraDocker/" + i.tfbin.GetEnvironmentFiles()[0])
	if err != nil {
		return nil, err
	}

	for _, container := range containerToDestroy {
		err := ihpenv.RemoveMySqlContainer(container.index)
		if err != nil {
			return nil, err
		}

		err = ihpcredsenv.RemoveMysqlRootPassword(container.index)
		if err != nil {
			return nil, err
		}

		err = WriteCustomEnv("./TerraDocker/terraform.tfvars.json", ihpenv)
		if err != nil {
			return nil, err
		}
		err = WriteCredsEnv("./TerraDocker/"+i.tfbin.GetEnvironmentFiles()[0], ihpcredsenv)
		if err != nil {
			return nil, err
		}
	}

	for _, container := range containerToCreate {
		ihpenv.AddMySqlContainer(container.GetUid(), 0, container.GetMySqlRootPassword())
		ihpcredsenv.AddMysqlRootPassword(container.GetMySqlRootPassword())

		err = WriteCustomEnv("./TerraDocker/terraform.tfvars.json", ihpenv)
		if err != nil {
			return nil, err
		}
		err = WriteCredsEnv("./TerraDocker/"+i.tfbin.GetEnvironmentFiles()[0], ihpcredsenv)
		if err != nil {
			return nil, err
		}
	}

	state, err := tfbin.RunTerra(tf, i.GetEnvironmentFiles()...)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (t *TerraformBinary) GetState() (*tfjson.State, error) {
	tf, err := t.GetInstance()
	if err != nil {
		return nil, err
	}
	state, err := tf.Show(context.Background())
	return state, err
}
