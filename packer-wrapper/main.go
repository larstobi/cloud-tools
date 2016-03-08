// passing command-line arguments to packer as-is
//
// To link the correct password to an environment, a wrapper.yml
// file has to be placed in the same directory as your *.tf files
//
//  Example wrapper.yml:
//
//  ---
//
//  secret-vars:
//
//    - name: CLOUDSTACK_API_KEY
//      key: CloudStack/accounts/admin/API_KEY
//
//    - name: CLOUDSTACK_SECRET_KEY
//      key: CloudStack/accounts/admin/SECRET_KEY
//
//  vars:
//
//    - name: CLOUDSTACK_API_URL
//      value: https://digipost-prod.cloudservices.no/client/api
//
package main

import (
	"os"

	"github.com/digipost/cloud-tools/config"
	"github.com/digipost/cloud-tools/wrapper"
)

// packer-wrapper will get secrets from your pass password store,
// setup an environment containing secrets and execute packer,
// passing command-line arguments to terraform as-is
func main() {

	config := config.ParseDefaultCloudConfig()
	secEnv := wrapper.GetEnvironmentVariablesForSecrets(config.SecretVariables[:])
	env := wrapper.GetEnvironmentVariablesForValues(config.Variables[:])
	wrapper.ExecuteCommand("packer", os.Args[1:], append(secEnv, env...))

}
