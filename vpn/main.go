package main

import (
	"fmt"
	"github.com/xanzy/go-cloudstack/cloudstack"
	"gopkg.in/gcfg.v1"
	"os"
	"strconv"
	"strings"
)

// Utility to enable remote access VPN on an 
// Apache Cloudstack VPC 
// 
// If 10.x.0.0/16 is the CIDR of your VPC,
// the vpn concentrator will distribute IPs in  
// Will give you an IP in the 10.(x+100).0.0 CIDR
// 
// To enable routing to this network:
//
// sudo route add 10.x.0.0/16 10.(x+100).0.1
// 
func main() {

	argsWithoutProg := os.Args[1:]

	var cfg config

	if nbArgs := len(argsWithoutProg); nbArgs != 2 {
		fmt.Printf("Enable remote VPN access on VPC\n")
		fmt.Printf("Usage: vpn <cloudstack.ini> <vpcname>\n")
		os.Exit(1)
	} else if err := gcfg.ReadFileInto(&cfg, argsWithoutProg[0]); err != nil {
		fmt.Printf("Error occured while reading config file: %s\n", err.Error())
		os.Exit(1)
	}

	apiurl := cfg.Cloudstack.Endpoint
	apikey := cfg.Cloudstack.Key
	secret := cfg.Cloudstack.Secret

	vpcName := argsWithoutProg[1]

	client := cloudstack.NewClient(apiurl, apikey, secret, true)
	asyncClient := cloudstack.NewAsyncClient(apiurl, apikey, secret, true)

	if vpcId, err := findVpcId(client, vpcName); err != nil {
		fmt.Printf("Failed to find id for VPC \"%s\": %s\n", vpcName, err.Error())
	} else {

		fmt.Printf("VPC id %s found for vpc name \"%s\"\n", vpcId, vpcName)

		if ipAddressId, err := findPublicIPAddressForVPC(client, vpcId); err != nil {
			fmt.Printf("Failed to find public IP address for VPC id %s: %s", vpcId, err.Error())
		} else {

			var vpn *cloudstack.RemoteAccessVpn

			if vpnExisting, err := findRemoteAccessVPN(client, ipAddressId); err != nil {
				fmt.Printf("Failed to find remote access VPN: %s", err.Error)
			} else if vpnExisting == nil {

				fmt.Printf("Remote Access VPN not enabled for VPC, creating new one\n")

				if vpnAddressRange, err := findVPNAddressRange(client, vpcId); err != nil {
					fmt.Printf("Failed to find cidr range: %s\n", err)
				} else {
					fmt.Printf("VPN address range %s \n", vpnAddressRange)

					if vpnCreated, err := createRemoteAccessVPN(asyncClient, ipAddressId, vpnAddressRange); err != nil {
						fmt.Printf("Failed to create new remote access VPN: %s", err.Error())
					} else {
						vpn = vpnCreated
					}

				}

			} else {
				vpn = vpnExisting
			}

			if vpn != nil {

				fmt.Printf("VPN connection details for VPC \"%s\":\n", vpcName)
				fmt.Printf("IP address: %s\n", vpn.Publicip)
				fmt.Printf("Preshared secret: %s\n", vpn.Presharedkey)

			}

		}

	}

}

func findVPNAddressRange(client *cloudstack.CloudStackClient, vpcId string) (string, error) {

	service := cloudstack.NewVPCService(client)
	params := service.NewListVPCsParams()
	params.SetId(vpcId)

	if vpcs, err := service.ListVPCs(params); err != nil {
		return "", err
	} else if vpcs.Count == 1 {
		cidr := vpcs.VPCs[0].Cidr
		fmt.Printf("CIDR range for VPC %s is %s\n", vpcId, cidr)

		return calculateVpnCidrRange(cidr), nil
	} else {
		return "", fmt.Errorf("VPC with id %s does not exist", vpcId)
	}

}

func calculateVpnCidrRange(vpcCidr string) string {
	address := strings.Split(vpcCidr, "/")[0]
	octets := strings.Split(address, ".")
	octet, _ := strconv.Atoi(octets[1])
	start := octets[0] + "." + strconv.Itoa(octet+100) + "." + octets[2] + ".1"
	end := octets[0] + "." + strconv.Itoa(octet+100) + "." + octets[2] + ".32"
	return start + "-" + end
}

func findRemoteAccessVPN(client *cloudstack.CloudStackClient, ipAddressId string) (*cloudstack.RemoteAccessVpn, error) {

	service := cloudstack.NewVPNService(client)
	params := service.NewListRemoteAccessVpnsParams()
	params.SetPublicipid(ipAddressId)

	if vpns, err := service.ListRemoteAccessVpns(params); err != nil {
		return nil, err
	} else if vpns.Count == 1 {
		return vpns.RemoteAccessVpns[0], nil
	} else {
		return nil, nil
	}

}

func createRemoteAccessVPN(client *cloudstack.CloudStackClient, ipAddressId string, addressRange string) (*cloudstack.RemoteAccessVpn, error) {

	service := cloudstack.NewVPNService(client)
	params := service.NewCreateRemoteAccessVpnParams(ipAddressId)
	params.SetFordisplay(true)
	params.SetOpenfirewall(true)
	params.SetIprange(addressRange)

	if vpn, err := service.CreateRemoteAccessVpn(params); err != nil {
		return nil, err
	} else {
		// Keeping this at a minimum
		return &cloudstack.RemoteAccessVpn{Publicip: vpn.Publicip, Presharedkey: vpn.Presharedkey}, nil
	}

}

func findPublicIPAddressForVPC(client *cloudstack.CloudStackClient, vpcId string) (string, error) {

	service := cloudstack.NewAddressService(client)
	params := service.NewListPublicIpAddressesParams()
	params.SetVpcid(vpcId)
	params.SetIssourcenat(true)

	if addresses, err := service.ListPublicIpAddresses(params); err != nil {
		return "", err
	} else if addresses.Count == 1 {
		return addresses.PublicIpAddresses[0].Id, nil
	} else {
		return "", fmt.Errorf("Virtual router source NAT ip address for vpcid %s not found", vpcId)
	}

}

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


type config struct {
	Cloudstack struct {
		Endpoint string
		Key      string
		Secret   string
	}
}
