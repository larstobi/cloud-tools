package main

import (
	"fmt"
	"github.com/digipost/cloud-tools/config"
	"github.com/xanzy/go-cloudstack/cloudstack"
	"os"
	"github.com/digipost/cloud-tools/cloudstackutils"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Restart VPC\n")
		fmt.Printf("Usage: cloudstack-vpc-restart <vpcname>\n")
		os.Exit(1)
	}

	vpcName := os.Args[1]

	apiurl, apikey, secret := config.CloudstackClientConfig()
	client := cloudstack.NewClient(apiurl, apikey, secret, true)

	vpcService := cloudstack.NewVPCService(client)

	if vpcId, err := cloudstackutils.FindVpcId(client, vpcName); err != nil {
		fmt.Printf("Failed to find id for VPC \"%s\": %s\n", vpcName, err.Error())
	} else {
		if _, err := vpcService.RestartVPC(vpcService.NewRestartVPCParams(vpcId)); err != nil {
			fmt.Printf("Failed to restart VPC \"%s\": %s\n", vpcName, err.Error())
		} else {
			fmt.Printf("Restarting VPC \"%s\"... ")
		}

	}

}