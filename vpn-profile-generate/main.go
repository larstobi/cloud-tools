package main
import (
	"text/template"
	"os"
	b64 "encoding/base64"
	"os/exec"
	"fmt"
	"strings"
)

type VpnSettings struct {
	ConsentText         string
	SharedSecret        string
	Username            string
	Password            string
	Hostname            string
	UserDefinedName     string // Name in menu
	PayloadDescription  string
	PayloadDisplayName  string
	PayloadIdentifier   string
	PayloadOrganization string
	PayloadUUIDInternal string
	PayloadUUID         string
}

// Will generate a VPN profile that can be installed/removed using
//
// profiles -I -F my.profile
// profiles -R -F my.profile
//
// Original XML was generated using "Apple Configurator 2" that can be
// installed from the app store
func main() {

	vpnSettings := VpnSettings{
		ConsentText: "Installing VPN connection to Signering Ops Test...",
		SharedSecret: "SharedSecret",
		Username: "myusername",
		Password: "mypassword",
		Hostname: "opstest.signering.posten.no",
		UserDefinedName:"Signering Ops Test (P)",
		PayloadDescription: "This will install a VPN connection to Signering Ops Test",
		PayloadDisplayName: "Profile for Signering Ops Test VPN",
		PayloadIdentifier: "no.posten.signering.opstest.vpn",
		PayloadOrganization: "Digipost",
		PayloadUUIDInternal: uuidgen(),
		PayloadUUID: uuidgen(),

	}

	funcMap := template.FuncMap{
		"base64": base64,
	}

	t := template.New("VPN Profile")
	t.Funcs(funcMap)
	t.Parse(vpnProfile)
	t.Execute(os.Stdout, vpnSettings)

}

func base64(data string) (string) {
	return b64.StdEncoding.EncodeToString([]byte(data))
}

func uuidgen() (string) {
	out, _ := exec.Command("uuidgen").Output()
	// get rid of trailing new line
	uuid := strings.TrimSpace(string(out))
	return fmt.Sprintf("%s", uuid)
}

const vpnProfile =
`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>ConsentText</key>
	<dict>
		<key>default</key>
		<string>{{ .ConsentText }}</string>
	</dict>
	<key>PayloadContent</key>
	<array>
		<dict>
			<key>IPSec</key>
			<dict>
				<key>AuthenticationMethod</key>
				<string>SharedSecret</string>
				<key>LocalIdentifierType</key>
				<string>KeyID</string>
				<key>SharedSecret</key>
				<data>
				{{ .SharedSecret | base64 }}
				</data>
			</dict>
			<key>IPv4</key>
			<dict>
				<key>OverridePrimary</key>
				<integer>0</integer>
			</dict>
			<key>PPP</key>
			<dict>
				<key>AuthName</key>
				<string>{{ .Username }}</string>
				<key>AuthPassword</key>
				<string>{{ .Password }}</string>
				<key>CommRemoteAddress</key>
				<string>{{ .Hostname }}</string>
			</dict>
			<key>PayloadDescription</key>
			<string>Configures VPN settings</string>
			<key>PayloadDisplayName</key>
			<string>VPN</string>
			<key>PayloadIdentifier</key>
			<string>{{ .PayloadIdentifier }}.internal</string>
			<key>PayloadType</key>
			<string>com.apple.vpn.managed</string>
			<key>PayloadUUID</key>
			<string>{{ .PayloadUUIDInternal }}</string>
			<key>PayloadVersion</key>
			<real>1</real>
			<key>Proxies</key>
			<dict>
				<key>HTTPEnable</key>
				<integer>0</integer>
				<key>HTTPSEnable</key>
				<integer>0</integer>
			</dict>
			<key>UserDefinedName</key>
			<string>{{ .UserDefinedName }}</string>
			<key>VPNType</key>
			<string>L2TP</string>
			<key>VendorConfig</key>
			<dict/>
		</dict>
	</array>
	<key>PayloadDescription</key>
	<string>{{ .PayloadDescription }}</string>
	<key>PayloadDisplayName</key>
	<string>{{ .PayloadDisplayName }}</string>
	<key>PayloadIdentifier</key>
	<string>{{ .PayloadIdentifier }}</string>
	<key>PayloadOrganization</key>
	<string>{{ .PayloadOrganization }}</string>
	<key>PayloadRemovalDisallowed</key>
	<false/>
	<key>PayloadType</key>
	<string>Configuration</string>
	<key>PayloadUUID</key>
	<string>{{ .PayloadUUID }}</string>
	<key>PayloadVersion</key>
	<integer>1</integer>
</dict>
</plist>`
