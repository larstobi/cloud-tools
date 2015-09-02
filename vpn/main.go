package main

import (
  "errors"
  "fmt"
  "github.com/xanzy/go-cloudstack/cloudstack"
)

func main() {

  url := "URL"
  key := "KEY"
  secret := "SECRET"

  vpcName := "prod"

  client := cloudstack.NewClient(url, key, secret, true)
  asyncClient := cloudstack.NewAsyncClient(url, key, secret, true)

  if vpcId, err := findVpcId(client, vpcName); err != nil {
    fmt.Printf(err.Error())
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

        if vpnCreated, err := createRemoteAccessVPN(asyncClient, ipAddressId); err != nil {
          fmt.Printf("Failed to create new remote access VPN: %s", err.Error())
        } else {
          vpn = vpnCreated
        }

      } else {
        vpn = vpnExisting
      }

      fmt.Printf("VPN connection details for VPC \"%s\":\n", vpcName)
      fmt.Printf("IP address: %s\n", vpn.Publicip)
      fmt.Printf("Preshared secret: %s\n", vpn.Presharedkey)

    }

  }

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

func createRemoteAccessVPN(client *cloudstack.CloudStackClient, ipAddressId string) (*cloudstack.RemoteAccessVpn, error) {

  service := cloudstack.NewVPNService(client)
  params := service.NewCreateRemoteAccessVpnParams(ipAddressId)
  params.SetFordisplay(true)
  params.SetOpenfirewall(true)

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
    return "", errors.New(fmt.Sprintf("Virtual router source NAT ip address for vpcid %s not found", vpcId))
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
    return "", errors.New(fmt.Sprintf("VPC %s not found", vpcName))
  }

}
