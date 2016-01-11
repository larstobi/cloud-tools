package main

import (
	"fmt"
	"github.com/digipost/cloud-tools/config"
	"github.com/xanzy/go-cloudstack/cloudstack"
	"sort"
	"os"
)

// Utility will return a list of templates sorted by date
// Use optional argument to filter using keyword
func main() {

	var keyword string  = ""

	if len(os.Args) == 2 {
		keyword = os.Args[1]
	}

	apiurl, apikey, secret := config.CloudstackClientConfig()
	client := cloudstack.NewClient(apiurl, apikey, secret, true)
	templateService := cloudstack.NewTemplateService(client)

	params := templateService.NewListTemplatesParams("community")
	if keyword != "" {
		params.SetKeyword(keyword)
	}

	if templates, err := templateService.ListTemplates(params); err != nil {
		panic(err)
	} else {
		sort.Sort(ByDate(templates.Templates))
		for _, template := range templates.Templates {
			fmt.Println(template.Created + " - " + template.Zonename + " - " + template.Name)
		}

	}

}

type ByDate []*cloudstack.Template

func (slice ByDate) Len() int {
	return len(slice)
}

func (slice ByDate) Less(i, j int) bool {
	return slice[i].Created < slice[j].Created
}

func (slice ByDate) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
