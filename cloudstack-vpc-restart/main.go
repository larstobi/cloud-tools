package main

import (
	"fmt"
	"github.com/digipost/cloud-tools/config"
	"github.com/xanzy/go-cloudstack/cloudstack"
	"os"
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

	if vpcId, err := findVpcId(client, vpcName); err != nil {
		fmt.Printf("Failed to find id for VPC \"%s\": %s\n", vpcName, err.Error())
	} else {
		if _, err := vpcService.RestartVPC(vpcService.NewRestartVPCParams(vpcId)); err != nil {
			fmt.Printf("Failed to restart VPC \"%s\": %s\n", vpcName, err.Error())
		} else {
			fmt.Printf("Restarting VPC \"%s\"... ")
		}

	}

}

// TODO extract this
func findVpcId(client *cloudstack.CloudStackClient, vpcName string) (string, error) {

	service := cloudstack.NewVPCService(client)
	params := service.NewListVPCsParams()
	params.SetName(vpcName)

	if vpcs, err := service.ListVPCs(params); err != nil {
		return "", err
	} else if vpcs.Count == 1 {
		return vpcs.VPCs[0].Id, nil
	} else {
		return "", fmt.Errorf("VPC %s does not exist", vpcName)
	}

}