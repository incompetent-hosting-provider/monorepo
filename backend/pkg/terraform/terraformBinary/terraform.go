package terraformbinary

import (
	"context"

	tfexec "github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	log "github.com/rs/zerolog/log"
)

func RunTerra(tf *tfexec.Terraform, varFiles ...string) (*tfjson.State, error) {
	// TODO: Add functionallity to include all the var files
	plan_diff, err := tf.Plan(context.Background(), tfexec.VarFile(varFiles[0]))
	if err != nil {
		log.Error().Msgf("Error planning terraform: %s", err)
		return nil, err
	}

	log.Debug().Msgf("Plan Diff: %v", plan_diff)

	if plan_diff {
		log.Info().Msg("Terraform plan has a difference between the current state and the plan")
	} else {
		log.Info().Msg("Terraform plan has no difference between the current state and the plan")
	}

	if plan_diff {
		// err := tf.Destroy(context.Background(), tfexec.VarFile("creds.tfvars"))
		// handleError(err, "Error destroying terraform")

		// TODO: use the varfiles variable
		err = tf.Apply(context.Background(), tfexec.VarFile("creds.tfvars.json"))
		if err != nil {
			log.Error().Msgf("Error applying terraform: %s", err)
			return nil, err
		}

		log.Info().Msg("Terraform apply complete")
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		log.Error().Msgf("Error showing terraform plan: %s", err)
		return nil, err
	}

	return state, nil
}
