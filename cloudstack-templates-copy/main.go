package main

import (
	"fmt"
	"github.com/digipost/cloud-tools/config"
	"github.com/xanzy/go-cloudstack/cloudstack"
	"os"
)

// Utility will copy templates matching keyword in from zone in to to zone
func main() {



	if len(os.Args) != 4 {
		fmt.Printf("Copy templates into other zone\n")
		fmt.Printf("Usage: cloudstack-templates-copy <keyword> <from zone> <to zone>\n")
		os.Exit(1)
	}

	keyword := os.Args[1]
	fromZone := os.Args[2]
	toZone := os.Args[3]

	apiurl, apikey, secret := config.CloudstackClientConfig()
	client := cloudstack.NewClient(apiurl, apikey, secret, true)


	// Lookup zone ids

	service := cloudstack.NewZoneService(client)

	var fromId = ""
	if from, _, err := service.GetZoneByName(fromZone); err != nil {
		panic(err)
	} else {
		fromId = from.Id
	}

	var toId = ""
	if to, _, err := service.GetZoneByName(toZone); err != nil {
		panic(err)
	} else {
		toId = to.Id
	}

	// Lookup templates

	templateService := cloudstack.NewTemplateService(client)

	// Can only own templates
	listTemplatesParams := templateService.NewListTemplatesParams("self")

	listTemplatesParams.SetKeyword(keyword)
	listTemplatesParams.SetZoneid(fromId)

	if templates, err := templateService.ListTemplates(listTemplatesParams); err != nil {
		panic(err)
	} else {
		for _, template := range templates.Templates {

			fmt.Println("Copy template", template.Name, "(", template.Id, ") from zone \"", fromZone, "\" (", fromId, ") to zone \"", toZone, "\" with Id (", toId, ")")

			copyTemplateParams := templateService.NewCopyTemplateParams(toId, template.Id)
			copyTemplateParams.SetSourcezoneid(fromId)
			if res, err := templateService.CopyTemplate(copyTemplateParams); err != nil {
				panic(err)
			} else {
				fmt.Println(res.Name)
			}

		}

	}


}