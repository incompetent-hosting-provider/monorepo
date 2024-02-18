package terraform

import (
	incompetenthostingprovider "incompetent-hosting-provider/backend/pkg/terraform/incompetentHostingProvider"
	"incompetent-hosting-provider/backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func init() {
	// This will install the terraform binary if not present
	incompetenthostingprovider.NewTfBin(tf_bin_dir, tf_version, tf_cwd_dir, tf_envs_names)
}

func SampleHandler(c *gin.Context) {
	tf_env, err := incompetenthostingprovider.ReadCustomEnv(tf_env_filepath)

	if err != nil {
		util.ThrowInternalServerErrorException(c, "Internal Server Error")
	}

	err = incompetenthostingprovider.WriteCustomEnv(tf_env_filepath, tf_env)
	if err != nil {
		util.ThrowInternalServerErrorException(c, "Internal Server Error")
	}

	ihpTfBin := incompetenthostingprovider.NewTfBin(tf_bin_dir, tf_version, tf_cwd_dir, tf_envs_names)

	state, err := ihpTfBin.GetState()

	log.Info().Msgf("%v", state)
}
