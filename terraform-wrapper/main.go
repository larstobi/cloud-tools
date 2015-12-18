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

/* Example wrapper.yml:

---

secret-vars:

  - name: AWS_ACCESS_KEY_ID
    key: Amazon/route53/ACCOUNT_ID

  - name: AWS_SECRET_ACCESS_KEY
    key: Amazon/route53/SECRET_KEY

vars:

  - name: AWS_DEFAULT_REGION
    value: eu-central-1

*/

type SecretVariable struct {
	Name string
	Key  string
}

type Variable struct {
	Name  string
	Value string
}

type Variables struct {
	SecretVariables []SecretVariable `secret-vars`
	Variables       []Variable       `vars`
}

func main() {

	parseWrapperConfig()
	argsWithoutProg := os.Args[1:]
	accountid := getPasswordFor("CloudStack/apikeys/ACCOUNT_ID")
	fmt.Println(accountid)
	executeTerraform(argsWithoutProg)

}

func parseWrapperConfig() {

	dir, _ := os.Getwd()
	filename, _ := filepath.Abs(dir + "/wrapper.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config Variables

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	// TODO remove
	fmt.Println("Secret var 0 : " + config.SecretVariables[0].Name)
	fmt.Println("Secret var 0 : " + config.SecretVariables[0].Key)
	fmt.Println("Var 0 key : " + config.Variables[0].Name)
	fmt.Println("Var 0 key : " + config.Variables[0].Value)

}



func executeTerraform(args []string) {
	cmd := exec.Command("/Users/landro/terraform_0.6.6_darwin_amd64/terraform", args...)

	var environment []string
	environment = append(environment,
		"AWS_ACCESS_KEY_ID=abc",
		"AWS_SECRET_ACCESS_KEY=eef",
		"AWS_DEFAULT_REGION=eu-abd")

	cmd.Env = environment
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	defer cmd.Wait()
}

func getPasswordFor(key string) string {

	cmd := exec.Command("/usr/local/bin/pass", key)
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
	var allText string

	for buff.Scan() {
		allText += buff.Text()
	}

	// TODO remove
	fmt.Printf("abc")
	fmt.Printf(allText)
	fmt.Printf("def")

	return allText

}
