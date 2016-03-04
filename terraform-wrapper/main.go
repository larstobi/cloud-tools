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
	start := time.Now()
	fmt.Println("Started terraform operation at:", start)
	config := config.ParseDefaultCloudConfig()
	secEnv := wrapper.GetEnvironmentVariablesForSecrets(config.SecretVariables[:])
	env := wrapper.GetEnvironmentVariablesForValues(config.Variables[:])
	wrapper.ExecuteCommand("terraform", os.Args[1:], append(secEnv, env...))
	stop := time.Now()
	fmt.Println("Started terraform operation at:", start)
	fmt.Println("Finished terraform operation at:", stop)
	duration := stop.Unix() - start.Unix()
	fmt.Println("Total duration (seconds): ", duration )
}

