package terraformbinary

import (
	"context"

	tfexec "github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

type TerraformBinaryInterface interface {
	NewTfBin(install_directory string, required_version string, working_directory string, environment_files []string) *TerraformBinary
	GetInstance() (*tfexec.Terraform, error)
	GetState() (*tfjson.State, error)
	GetEnvironmentFiles() []string
}

type TerraformBinary struct {
	install_directory string
	required_version  string
	working_directory string
	environment_files []string
}

func NewTfBin(install_directory string, required_version string, working_directory string, environment_files []string) *TerraformBinary {
	return &TerraformBinary{install_directory, required_version, working_directory, environment_files}
}

func (t *TerraformBinary) GetInstance() (*tfexec.Terraform, error) {
	tf, err := getInstance(t.install_directory, t.required_version, t.working_directory)
	return tf.asTerraform(), err
}

func (t *TerraformBinary) GetState() (*tfjson.State, error) {
	tf, err := t.GetInstance()
	if err != nil {
		return nil, err
	}
	state, err := tf.Show(context.Background())
	return state, err
}

func (t *TerraformBinary) GetEnvironmentFiles() []string {
	return t.environment_files
}
