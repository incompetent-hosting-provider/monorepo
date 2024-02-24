package terraformbinary

import (
	"context"

	tfexec "github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	log "github.com/rs/zerolog/log"
)

func RunTerra(tf *tfexec.Terraform, tfFilesPath string, varFiles ...string) (*tfjson.State, error) {
	tfPlanOpts := []tfexec.PlanOption{}
	log.Debug().Msgf("Terraform files path: %s", tfFilesPath)
	for _, varFile := range varFiles {
		varFileOpt := tfexec.VarFile(varFile)
		tfPlanOpts = append(tfPlanOpts, varFileOpt)
	}

	plan_diff, err := tf.Plan(context.TODO(), tfPlanOpts...)
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
		// TODO: Rethink if destroying beforehand is necessary
		// err := tf.Destroy(context.TODO(), tfexec.VarFile("creds.tfvars"))
		// handleError(err, "Error destroying terraform")

		tfApplyOpts := []tfexec.ApplyOption{}
		for _, varFile := range varFiles {
			varFileOpt := tfexec.VarFile(varFile)
			tfApplyOpts = append(tfApplyOpts, varFileOpt)
		}

		err = tf.Apply(context.TODO(), tfApplyOpts...)
		if err != nil {
			log.Error().Msgf("Error applying terraform: %s", err)
			return nil, err
		}

		log.Info().Msg("Terraform apply complete")
	}

	state, err := tf.Show(context.TODO())
	if err != nil {
		log.Error().Msgf("Error showing terraform plan: %s", err)
		return nil, err
	}

	return state, nil
}
