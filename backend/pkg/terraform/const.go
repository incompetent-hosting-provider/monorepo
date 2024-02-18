package terraform

import "os"

var cwd, _ = os.Getwd()
var tf_bin_dir string = cwd + "/bins"

const tf_version string = "1.7.1"

var tf_cwd_dir string = cwd + "/terraform"

var tf_env_filepath = tf_cwd_dir + "/terraform.tfvars.json"

var tf_envs_names []string = []string{
	"creds.tfvars.json",
}
