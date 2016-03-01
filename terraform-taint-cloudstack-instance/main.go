package main

import (
	"fmt"
	"github.com/digipost/cloud-tools/terraform"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Taint cloudstack instance safely\n")
		fmt.Printf("Will remove nic\n")
		fmt.Printf("Usage:cloudstack-taint-cloudstack-instance <instance resource id>\n")
		fmt.Printf("Example: cloudstack-taint-cloudstack-instance cloudstack_instance.db.1\n")
		os.Exit(1)
	}

	cloudstackInstanceId := os.Args[1]

	// Read Terraform state
	if state, err := terraform.ReadState("terraform.tfstate"); err != nil {
		fmt.Printf("Unable to read terraform.tfstate: %s", err.Error())
		os.Exit(1)
	} else {

		if resourceState, ok := state.Modules[0].Resources[cloudstackInstanceId]; ok {

			vmName := resourceState.Primary.Attributes["name"]
			fmt.Println("terraform taint " + cloudstackInstanceId)

			for resourceId, resource := range state.Modules[0].Resources {
				if resource.Type == "cloudstack_nic" && resource.Primary.Attributes["virtual_machine"] == vmName {
					fmt.Println("terraform taint " + resourceId)
				}
			}

		} else {
			fmt.Println("Instance id not found")
			os.Exit(1)
		}
	}

}
