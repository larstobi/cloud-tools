package config

import (
//"github.com/digipost/cloud-tools/config"
)

// Returns CloudStack relevant client config
// apiurl, apikey and secret retrieved from the pass
// password store
func CloudstackClientConfig() (string, string, string) {

	var apiurl string
	var apikey string
	var secret string

	cloudConfig := ParseDefaultCloudConfig()
	secVars := cloudConfig.SecretVariables
	for i := 0; i < len(secVars); i++ {

		if secVars[i].Name == "CLOUDSTACK_API_KEY" {
			apikey = GetPasswordFor(secVars[i].Key)
		}

		if secVars[i].Name == "CLOUDSTACK_SECRET_KEY" {
			secret = GetPasswordFor(secVars[i].Key)
		}

	}
	vars := cloudConfig.Variables
	for i := 0; i < len(vars); i++ {
		if vars[i].Name == "CLOUDSTACK_API_URL" {
			apiurl = vars[i].Value
		}
	}

	return apiurl, apikey, secret

}
