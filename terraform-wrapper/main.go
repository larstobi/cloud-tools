package main

import (
	"github.com/digipost/cloud-tools/config"
	"github.com/digipost/cloud-tools/wrapper"
	"os"
	"fmt"
	"time"
)

// terraform-wrapper will get secrets from your pass password store,
// setup an environment containing secrets and execute terraform,
// passing command-line arguments to terraform as-is
func main() {
	fmt.Println("Starting terraform operation at:", time.Now())
	config := config.ParseDefaultCloudConfig()
	secEnv := wrapper.GetEnvironmentVariablesForSecrets(config.SecretVariables[:])
	env := wrapper.GetEnvironmentVariablesForValues(config.Variables[:])
	wrapper.ExecuteTerraform("terraform", os.Args[1:], append(secEnv, env...))
	fmt.Println("Ending terraform operation at:", time.Now())
}

