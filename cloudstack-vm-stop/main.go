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
		fmt.Printf("Stop VMs in VPC\n")
		fmt.Printf("Usage: cloudstack-vm-stop <vpcname>\n")
		os.Exit(1)
	}

	vpcName := os.Args[1]

	apiurl, apikey, secret := config.CloudstackClientConfig()
	client := cloudstack.NewClient(apiurl, apikey, secret, true)

	if vpcId, vpcName, err := cloudstackutils.FindVpcId(client, vpcName); err != nil {
		fmt.Printf("Failed to find id for VPC \"%s\": %s\n", vpcName, err.Error())
	} else {

		vmService := cloudstack.NewVirtualMachineService(client)
		virtualMachinesParams := vmService.NewListVirtualMachinesParams()
		virtualMachinesParams.SetVpcid(vpcId)
		virtualMachinesParams.SetState("Running")

		if vms, err := vmService.ListVirtualMachines(virtualMachinesParams); err != nil {
			fmt.Printf("Failed to list VMs in VPC \"%s\": %s\n", vpcName, err.Error())
		} else {

			if vms.Count > 0 {

				for _, vm := range vms.VirtualMachines {
					fmt.Printf("Stopping VM %s / %s \n", vm.Displayname, vm.Name)
					if _, err := vmService.StopVirtualMachine(vmService.NewStopVirtualMachineParams(vm.Id)); err != nil {
						fmt.Printf("Failed to stop VM with id %s: %s\n", vm.Id, err.Error())
					} else {
						fmt.Printf("VM with id %s stopped\n", vm.Id)
					}
				}
			} else {
				fmt.Printf("No running VMs in %s\n", vpcName)
			}
		}

	}

}