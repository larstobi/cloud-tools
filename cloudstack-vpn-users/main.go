package main

import (
	"io/ioutil"
	"github.com/digipost/cloud-tools/config"
	"github.com/xanzy/go-cloudstack/cloudstack"
	"fmt"
	"os"
)

// Utility will delete all existing vpn users and create new ones and store generated passwords
// in provided password store
func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Generate VPN users and save passwords in given pass password store\n")
		fmt.Printf("Usage: cloudstack-vpn-users <PASSWORD_STORE_DIR>\n")
		os.Exit(1)
	}

	passwordStoreDir := os.Args[1]

	// Find usernames
	fileinfos, _ := ioutil.ReadDir(fmt.Sprintf("%s/users", passwordStoreDir))
	users := make([]string, 0)
	for _, fileinfo := range fileinfos {
		if fileinfo.IsDir() {
			users = append(users, fileinfo.Name())
		}
	}

	// Generate passwords
	passwords := make([]string, 0)
	for _, user := range users {
		passwords = append(passwords, config.GeneratePasswordFor(passwordStoreDir, fmt.Sprintf("users/%s/vpn", user), 32))
	}

	apiurl, apikey, secret := config.CloudstackClientConfig()

	client := cloudstack.NewAsyncClient(apiurl, apikey, secret, true)

	//	domains := getDomains(client)
	//	for _, domain := range domains {
	//
	//		fmt.Println(domain.Name + "/" + domain.Path + "(" + domain.Id + ")")
	//
	//
	//	}


	// List all existing vpn accounts
	// TODO Fix this

	//(ROOT/digipost/sign)
	digipost_sign_domain := "191b7185-0d65-47a4-b4c5-22cadef03baf"

	// Get all accounts in domain
	accounts := getAccounts(client, digipost_sign_domain)

	// Delete all existing VPN users from all accounts in domain
	for _, account := range accounts {
		removeAllVpnUsers(client, digipost_sign_domain, account.Name)
	}

	// Add new remote access VPN users to all accounts in domain
	for _, account := range accounts {
		fmt.Printf("--------------------------------------------------------------------------\n")
		fmt.Printf("Adding remote access VPN users for account: %s\n", account.Name)
		fmt.Printf("--------------------------------------------------------------------------\n")
		for index, user := range users {
			addVpnUser(client, user, passwords[index], digipost_sign_domain, account.Name)
		}
	}

}

func removeAllVpnUsers(client *cloudstack.CloudStackClient, domainid string, accountName string) {
	fmt.Printf("--------------------------------------------------------------------------\n")
	fmt.Printf("Removing remote access VPN users for account: %s\n", accountName)
	fmt.Printf("--------------------------------------------------------------------------\n")
	vpnService := cloudstack.NewVPNService(client)
	vpnUsers := getVpnUsers(client, domainid, accountName)
	for _, vpnUser := range vpnUsers {
		params := vpnService.NewRemoveVpnUserParams(vpnUser.Username)
		params.SetAccount(accountName)
		params.SetDomainid(domainid)
		if _, err := vpnService.RemoveVpnUser(params); err != nil {
			fmt.Printf("Failed to remove remote access VPN for user %s: %s\n", vpnUser.Username, err.Error())
		} else {
			fmt.Printf("Removed remote access VPN for user %s (account: %s, domainid: %s)\n", vpnUser.Username, accountName, domainid)
		}
	}
}

func addVpnUser(client *cloudstack.CloudStackClient, username string, password string, domainid string, accountName string) {

	vpnService := cloudstack.NewVPNService(client)

	params := vpnService.NewAddVpnUserParams(password, username)
	params.SetAccount(accountName)
	params.SetDomainid(domainid)
	if _, err := vpnService.AddVpnUser(params); err != nil {
		fmt.Printf("Failed to create new remote access VPN: %s\n", err.Error())
	} else {
		fmt.Printf("Added user %s to VPN (account: %s, domainid: %s)\n", username, accountName, domainid)
	}

}

func getAccounts(client *cloudstack.CloudStackClient, domainid string) ([]*cloudstack.Account) {
	accountService := cloudstack.NewAccountService(client)
	params := accountService.NewListAccountsParams()
	params.SetDomainid(domainid)
	params.SetListall(true)
	accounts, _ := accountService.ListAccounts(params)
	return accounts.Accounts
}

func getVpnUsers(client *cloudstack.CloudStackClient, domainid string, account string) ([]*cloudstack.VpnUser) {
	vpnService := cloudstack.NewVPNService(client)
	params := vpnService.NewListVpnUsersParams()
	params.SetAccount(account)
	params.SetDomainid(domainid)
	params.SetListall(true)
	users, _ := vpnService.ListVpnUsers(params)
	return users.VpnUsers
}

func getDomains(client *cloudstack.CloudStackClient) ([]*cloudstack.Domain) {
	domainService := cloudstack.NewDomainService(client)
	params := domainService.NewListDomainsParams()
	params.SetListall(true)
	domains, _ := domainService.ListDomains(params)
	return domains.Domains
}
