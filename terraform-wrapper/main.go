package main

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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

	config := parseWrapperConfig()
	secEnv := getEnvironmentVariablesForSecrets(config.SecretVariables[:])
	env := getEnvironmentVariablesForValues(config.Variables[:])
	executeTerraform(os.Args[1:], append(secEnv, env...))

}

func getEnvironmentVariablesForSecrets(secretVars []SecretVariable) []string {
	var environment []string
	for i := 0; i < len(secretVars); i++ {
		environment = append(environment, secretVars[i].Name+"="+getPasswordFor(secretVars[i].Key))
	}
	return environment
}

func getEnvironmentVariablesForValues(vars []Variable) []string {
	var environment []string
	for i := 0; i < len(vars); i++ {
		environment = append(environment, vars[i].Name+"="+vars[i].Value)
	}
	return environment
}

func parseWrapperConfig() Config {

	dir, _ := os.Getwd()
	filename, _ := filepath.Abs(fmt.Sprintf("%s%c%s", dir, os.PathSeparator, "wrapper.yml"))
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	return config

}

func executeTerraform(args []string, environment []string) {

	cmd := exec.Command("terraform", args...)

	cmd.Env = append(environment, "PATH="+os.Getenv("PATH"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	defer cmd.Wait()
}

func getPasswordFor(key string) string {

	cmd := exec.Command("pass", key)
	//cmd.Env = ()
	// Ask for gpg password if necessary
	cmd.Stdin = os.Stdin

	// capture the output and error pipes
	stdout, _ := cmd.StdoutPipe()
	//stderr, _ := cmd.StderrPipe()

	cmd.Start()

	// Don't let main() exit before our command has finished running
	// doesn't block
	defer cmd.Wait()

	buff := bufio.NewScanner(stdout)
	var password string

	for buff.Scan() {
		password += buff.Text()
	}

	return password

}

type SecretVariable struct {
	Name string
	Key  string
}

type Variable struct {
	Name  string
	Value string
}

type Config struct {
	SecretVariables []SecretVariable `secret-vars`
	Variables       []Variable       `vars`
}
