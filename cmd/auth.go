package cmd

import (
	"fmt"
	"os"

	"github.com/isaaacqinh/haos-ldap-auth/internal/gateway"
)

func Auth() {
	config := gateway.Config{}
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

	res, err := gateway.SearchUser(adConn, config, userCreds.Username)
	if err != nil {
		panic(err)
	}

	authString = fmt.Sprintf("name = %s\n", res.Username)
	authString += fmt.Sprintf("group = %s\n", res.Group)

	err = gateway.TryBind(adConn, config, res.Login, userCreds.Password)
	if err != nil {
		os.Exit(1)
	}

	fmt.Print(authString)
	os.Exit(0)
}
