package main

import (
	"fmt"
	"os"
	"github.com/digipost/cloud-tools/config"
	"github.com/xanzy/go-cloudstack/cloudstack"
	"github.com/digipost/cloud-tools/terraform"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Stop VM disk is attached to\n")
		fmt.Printf("Usage:cloudstack-stop-vm-for-disk <disk resource id> [-f] \n")
		fmt.Printf("Example: cloudstack-stop-vm-for-disk cloudstack_disk.db.1\n")
		os.Exit(1)
	}

	diskResourceId := os.Args[1]

	// Should we force stop
	forced := false
	if (len(os.Args) == 3) {

		if os.Args[2] == "-f" {
			fmt.Printf("Will try to forcefully stop vm\n")
			forced = true
		} else {
			fmt.Printf("Unrecognized second parameter. Should be -f\n")
			os.Exit(1)
		}

	}

	// Read Terraform state
	if state, err := terraform.ReadState("terraform.tfstate"); err != nil {
		fmt.Printf("Unable to read terraform.tfstate: %s", err.Error())
		os.Exit(1)
	} else {
		if resourceState, ok := state.Modules[0].Resources[diskResourceId]; ok {

			vmId := resourceState.Primary.Attributes["virtual_machine"]
			apiurl, apikey, secret := config.CloudstackClientConfig()
			client := cloudstack.NewAsyncClient(apiurl, apikey, secret, true)

			vmService := cloudstack.NewVirtualMachineService(client)
			params := vmService.NewStopVirtualMachineParams(vmId)
			params.SetForced(forced)

			if res, err := vmService.StopVirtualMachine(params); err != nil {
				fmt.Printf("Unable to stop vm: %s", err.Error())
				os.Exit(1)
			} else {
				fmt.Printf("Stopped vm with id %s\n", vmId)
				fmt.Printf("State is %s\n", res.State)
				os.Exit(0)
			}
		} else {
			fmt.Println("Disk id not found")
			os.Exit(1)
		}

	}

}
