package main

import (
	"github.com/digipost/cloud-tools/config"
	"os"
	"os/exec"
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
	secEnv := getEnvironmentVariablesForSecrets(config.SecretVariables[:])
	env := getEnvironmentVariablesForValues(config.Variables[:])
	executeTerraform(os.Args[1:], append(secEnv, env...))

}

func getEnvironmentVariablesForSecrets(secretVars []config.SecretVariable) []string {
	var environment []string
	for i := 0; i < len(secretVars); i++ {
		environment = append(environment, secretVars[i].Name+"="+config.GetPasswordFor(secretVars[i].Key))
	}
	return environment
}

func getEnvironmentVariablesForValues(vars []config.Variable) []string {
	var environment []string
	for i := 0; i < len(vars); i++ {
		environment = append(environment, vars[i].Name+"="+vars[i].Value)
	}
	return environment
}

func executeTerraform(args []string, environment []string) {

	cmd := exec.Command("terraform", args...)

	cmd.Env = append(environment, "PATH="+os.Getenv("PATH"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	defer cmd.Wait()
}
