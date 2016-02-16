package main
import (
	"github.com/digipost/cloud-tools/profile/vpn"
	"os"
	"fmt"
	"strings"
	"github.com/digipost/cloud-tools/config"
	"os/exec"
)


// Will generate a VPN profile that can be installed/removed using
//
// profiles -I -F my.profile
// profiles -R -F my.profile
// profiles -R -p <payloadidentifier>

// List profiles
// profiles -L
//
// Remove profiles
// profiles -R -p no.posten.x.x.x
//
func main() {

	if len(os.Args) != 4 {
		fmt.Println("Install VPN profiles")
		fmt.Println("Usage: vpn-profile-install <VPC-NAME> <VPN-USERNAME> <PASSWORD_STORE_DIR>")
		os.Exit(1)
	}

	vpcName := os.Args[1]
	username := os.Args[2]
	passwordStoreDir := os.Args[3]

	hostname := getVpnHostnameForVpcName(vpcName)

	passwordKey := fmt.Sprintf("users/%s/vpn", username)

	password := config.GetPasswordFromPasswordStoreFor(passwordKey, passwordStoreDir)
	sharedSecret := config.GetPasswordFromPasswordStoreFor(vpcName, passwordStoreDir)


	vpnSettings := &vpn.Settings{
		ConsentText: "Installing VPN connection...",
		SharedSecret: sharedSecret,
		Username: username,
		Password: password,
		Hostname: hostname,
		UserDefinedName: vpcName,
		PayloadDescription: fmt.Sprintf("This will install a VPN connection to %s", vpcName),
		PayloadDisplayName: fmt.Sprintf("Profile for %s", vpcName),
		PayloadIdentifier:reverse(hostname),
		PayloadOrganization: vpcName,
	}


	cmd := exec.Command("profiles", "-I", "-F", "-")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	writer, _ := cmd.StdinPipe()
	cmd.Start()
	defer cmd.Wait()

	vpn.GenerateVpnProfile(writer, vpnSettings)
	defer writer.Close()
}

func getVpnHostnameForVpcName(vpcName string) (string) {

	parts := strings.Split(vpcName, "-")
	zone := parts[1]
	systemEnv := strings.Split(parts[0], "_")
	systemShort := systemEnv[0]
	env := systemEnv[1]

	var domain string
	switch {
	case "sign" <= systemShort:
		domain = "signering.posten.no"
	case "mf" <= systemShort:
		domain = "meldingsformidler.digipost.no"
	case "dp" <= systemShort:
		domain = "digipost.no"
	}

	return fmt.Sprintf("vpn.%s.%s.%s", zone, env, domain)

}

func reverse(hostname string) (string) {
	result := make([]string, 0)
	parts := strings.Split(hostname, ".")
	numberOfSplits := len(parts)
	for i := 1; i <= numberOfSplits; i++ {
		result = append(result, parts[numberOfSplits - i])
	}
	return strings.Join(result, ".")
}