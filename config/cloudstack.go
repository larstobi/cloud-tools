package config

// Returns CloudStack relevant client config
// apiurl, apikey and secret retrieved from the pass
// password store
func CloudstackClientConfig() (string, string, string) {

	var apiurl string
	var apikey string
	var secret string

	cloudConfig := ParseDefaultCloudConfig()

	for _, secVar := range cloudConfig.SecretVariables {

		if secVar.Name == "CLOUDSTACK_API_KEY" {
			apikey = GetPasswordFor(secVar.Key)
		}

		if secVar.Name == "CLOUDSTACK_SECRET_KEY" {
			secret = GetPasswordFor(secVar.Key)
		}

	}

	for _, variable := range cloudConfig.Variables {
		if variable.Name == "CLOUDSTACK_API_URL" {
			apiurl = variable.Value
		}
	}

	return apiurl, apikey, secret

}
