package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"fmt"
	"strings"
	"github.com/digipost/cloud-tools/config"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Initialise VPN password store in given pass password store\n")
		fmt.Printf("Usage: pass-vpn-password-store-init <PASSWORD_STORE_DIR>\n")
		os.Exit(1)
	}

	pwdStore := os.Args[1]

	dir, _ := os.Getwd()
	absFilename, _ := filepath.Abs(fmt.Sprintf("%s%call.yml", dir, os.PathSeparator))
	yamlFile, err := ioutil.ReadFile(absFilename)

	if err != nil {
		panic(err)
	}

	var developerConfig DeveloperConfig

	err = yaml.Unmarshal(yamlFile, &developerConfig)
	if err != nil {
		panic(err)
	}

	// Initialise password store for all users
	validDevelopers := make([]string, 0)
	for _, developer := range developerConfig.Developers {
		if developer.isValid() {
			fmt.Println("pass init -p users/" + developer.VpnUser + " " + developer.GpgKeyId)
			config.InitialisePasswordStore(pwdStore, "users/" + developer.VpnUser, developer.GpgKeyId)
			validDevelopers = append(validDevelopers, developer.GpgKeyId)
		}
	}
	fmt.Println("pass init -p shared " + strings.Join(validDevelopers, " "))
	config.InitialisePasswordStore(pwdStore, "shared", validDevelopers...)

}

type DeveloperConfig struct {
	Developers []Developer
}

type Developer struct {
	GpgKeyId string  `gpg-key-id`
	VpnUser  string `vpn-user`
}

func (d *Developer) isValid() bool {
	return d.GpgKeyId != "" && d.VpnUser != ""
}