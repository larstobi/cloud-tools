package main

import "testing"

func TestHostnameIsResolveCorrectly(t *testing.T) {

	hostname := getVpnHostnameForVpcName("sign_opstest")

	if hostname != "vpn.opstest.signering.posten.no" {
		t.Error("Expected vpn.opstest.signering.posten.no, got ", hostname)
	}

}

