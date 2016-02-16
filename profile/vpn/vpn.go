package vpn

import (
	"text/template"
	b64 "encoding/base64"
	"os/exec"
	"fmt"
	"strings"
	"time"
	"io"
)

type Settings struct {
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
}

func (u Settings) PayloadUUID() string {
	return uuidgen();
}

func GenerateVpnProfile(writer io.Writer, vpnSettings *Settings) {

	funcMap := template.FuncMap{
		"base64": base64,
	}

	t := template.New("VPN Profile")
	t.Funcs(funcMap)
	t.Parse(vpnProfile)
	t.Execute(writer, vpnSettings)

}

func base64(data string) (string) {
	return b64.StdEncoding.EncodeToString([]byte(data))
}

func version() (int64) {
	return time.Now().Unix()
}

func uuidgen() (string) {
	out, _ := exec.Command("uuidgen").Output()
	// get rid of trailing new line
	uuid := strings.TrimSpace(string(out))
	return fmt.Sprintf("%s", uuid)
}

// Original XML was generated using "Apple Configurator 2" that can be
// installed from the app store
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
			<string>{{ .PayloadUUID }}</string>
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
