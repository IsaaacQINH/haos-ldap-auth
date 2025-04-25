package cmd

import (
	"fmt"
	"os"

	"github.com/isaaacqinh/haos-ldap-auth/internal/gateway"
	"github.com/isaaacqinh/haos-ldap-auth/internal/logger"
)

func Auth() {
	logger.WriteLog("auth.log", "==== Starting authentication ====")
	config := gateway.Config{}
	userCreds, err := gateway.GetEnv()
	authString := ""

	if err != nil {
		logger.WriteLog("auth.log", fmt.Sprintf("Error getting environment variables: %v", err))
		panic(err)
	}

	config.GetConf()

	adConn, err := gateway.ConnectAndBind(config)
	if err != nil {
		logger.WriteLog(config.Verbose.File, fmt.Sprintf("Error connecting to LDAP server: %v", err))
		panic(err)
	}
	defer adConn.Close()

	res, err := gateway.SearchUser(adConn, config, userCreds.Username)
	if err != nil {
		logger.WriteLog(config.Verbose.File, fmt.Sprintf("Error searching for user: %v", err))
		panic(err)
	}

	authString = fmt.Sprintf("name = %s\n", res.Username)
	authString += fmt.Sprintf("group = %s\n", res.Group)

	err = gateway.TryBind(adConn, config, res.Login, userCreds.Password)
	if err != nil {
		logger.WriteLog(config.Verbose.File, fmt.Sprintf("Error binding user: %v", err))
		os.Exit(1)
	}

	logger.WriteLog(config.Verbose.File, fmt.Sprintf("User %s authenticated successfully", res.Username))
	fmt.Print(authString)
	os.Exit(0)
}
