package main

import "testing"

func TestHostnameIsResolveCorrectly(t *testing.T) {

	hostname := getVpnHostnameForVpcName("sign_opstest")

	expected := "vpn.opstest.signering.posten.no"
	if hostname != expected {
		t.Error("Expected ", expected, ", got ", hostname)
	}

}

