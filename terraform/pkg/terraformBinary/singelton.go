package terraformbinary

import (
	"reflect"
	"sync"

	tfexec "github.com/hashicorp/terraform-exec/tfexec"
	log "github.com/rs/zerolog/log"
)

var lock = &sync.Mutex{}

type singleTfExec struct {
	tf_instance *tfexec.Terraform
}

var singleInstance *singleTfExec

func getInstance(binary_dir string, terraform_version string, terraform_dir string) (*singleTfExec, error) {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			log.Debug().Msgf("Creating %s instance now.", reflect.TypeOf(singleTfExec{}))
			terraform_executable, err := ensureTerraform(binary_dir, terraform_version, terraform_dir)
			if err != nil {
				log.Error().Msgf("Error ensuring terraform: %s", err)
				return nil, err
			}

			singleInstance = &singleTfExec{tf_instance: terraform_executable}
		} else {
			log.Debug().Msgf("%s instance already created. Returning existing instance.", reflect.TypeOf(singleTfExec{}))
		}
	} else {
		log.Debug().Msgf("%s instance already created. Returning existing instance.", reflect.TypeOf(singleTfExec{}))
	}

	return singleInstance, nil
}

// Allow conversion to tfexec.Terraform with a type assertion to *tfexec.Terraform
func (s *singleTfExec) asTerraform() *tfexec.Terraform {
	return s.tf_instance
}
