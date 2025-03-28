package cmd

import (
	"fmt"
	"os"

	"github.com/isaaacqinh/haos-ldap-auth/internal/gateway"
)

func Auth() {

	config := gateway.Config{}
	searchGroups := config.Groups[:0]
	userCreds, err := gateway.GetEnv()
	authString := ""

	if err != nil {
		panic(err)
	}

	config.GetConf()

	adConn, err := gateway.ConnectAndBind(config)
	if err != nil {
		panic(err)
	}
	defer adConn.Close()

	groups, err := gateway.GetGroups(adConn, config)
	if err != nil {
		panic(err)
	}

	for _, group := range groups {
		searchGroups = append(searchGroups, group.DN)
	}

	res, err := gateway.SearchUser(adConn, config, searchGroups, userCreds.Username)
	if err != nil {
		panic(err)
	}

	authString = fmt.Sprintf("name = %s\n", res.GetAttributeValue("displayName"))
	authString += fmt.Sprintf("group = %s\n", gateway.IsAdmin(res, config))

	err = gateway.TryBind(adConn, config, res.DN, userCreds.Password)
	if err != nil {
		os.Exit(1)
	}

	fmt.Print(authString)
	os.Exit(0)
}
