package main

import "os"
import (
	"github.com/hashicorp/terraform/terraform"
	"fmt"
	"github.com/digipost/cloud-tools/config"
	"github.com/xanzy/go-cloudstack/cloudstack"
)

func readState(path string) (*terraform.State, error) {
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		defer f.Close()
		state, _ := terraform.ReadState(f)
		return state, nil
	}
}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Start VM disk is attached to\n")
		fmt.Printf("Usage:cloudstack-start-vm-for-disk <disk resource id>\n")
		fmt.Printf("Example: cloudstack-start-vm-for-disk cloudstack_disk.db.1\n")
		os.Exit(1)
	}

	diskResourceId := os.Args[1]

	// Read Terraform state
	if state, err := readState("terraform.tfstate"); err != nil {
		fmt.Printf("Unable to read terraform.tfstate: %s", err.Error())
		os.Exit(1)
	} else {
		primary := state.Modules[0].Resources[diskResourceId].Primary
		vmId := primary.Attributes["virtual_machine"]

		apiurl, apikey, secret := config.CloudstackClientConfig()
		client := cloudstack.NewAsyncClient(apiurl, apikey, secret, true)

		vmService := cloudstack.NewVirtualMachineService(client)
		params := vmService.NewStartVirtualMachineParams(vmId)

		if res, err := vmService.StartVirtualMachine(params); err != nil {
			fmt.Printf("Unable to start vm: %s", err.Error())
			os.Exit(1)
		} else {
			fmt.Printf("Started vm with id %s\n", vmId)
			fmt.Printf("State is %s\n", res.State)
			os.Exit(0)
		}

	}

}
