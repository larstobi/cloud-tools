package main

import (
	"github.com/digipost/cloud-tools/config"
  "github.com/digipost/cloud-tools/wrapper"
	"os"
)

// terraform-wrapper will get secrets from your pass password store,
// setup an environment containing secrets and execute terraform,
// passing command-line arguments to terraform as-is
//
// To link the correct password to an environment, a wrapper.yml
// file has to be placed in the same directory as your *.tf files
//
//	Example wrapper.yml:
//
//	---
//
//	secret-vars:
//
//	  - name: AWS_ACCESS_KEY_ID
//	    key: Amazon/route53/ACCOUNT_ID
//
//	  - name: AWS_SECRET_ACCESS_KEY
//	    key: Amazon/route53/SECRET_KEY
//
//	vars:
//
//	  - name: AWS_DEFAULT_REGION
//	    value: eu-central-1
//
func main() {

	config := config.ParseDefaultCloudConfig()
	secEnv := wrapper.GetEnvironmentVariablesForSecrets(config.SecretVariables[:])
	env := wrapper.GetEnvironmentVariablesForValues(config.Variables[:])
	wrapper.ExecuteTerraform("terraform", os.Args[1:], append(secEnv, env...))

}

