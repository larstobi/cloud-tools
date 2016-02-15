package cloudstackutils

import (
	"github.com/xanzy/go-cloudstack/cloudstack"
	"fmt"
)

func FindVpcId(client *cloudstack.CloudStackClient, vpcName string) (string, string, error) {

	service := cloudstack.NewVPCService(client)
	params := service.NewListVPCsParams()
	params.SetName(vpcName)

	if vpcs, err := service.ListVPCs(params); err != nil {
		return "", "", err
	} else if vpcs.Count == 1 {
		return vpcs.VPCs[0].Id, vpcs.VPCs[0].Name, nil
	} else {
		return "", "", fmt.Errorf("VPC \"%s\" does not exist, or too many VPCs matching name/filter", vpcName)
	}

}
