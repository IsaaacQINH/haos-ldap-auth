package gateway

import (
	"crypto/tls"
	"fmt"
	"slices"

	"github.com/go-ldap/ldap/v3"
)

type User struct {
	Login    string
	Username string
	Group    string
}

func ConnectAndBind(cfg Config) (*ldap.Conn, error) {
	var prot string
	var port int

	if cfg.TLS {
		prot = "ldaps"
		port = 636
	} else {
		prot = "ldap"
		port = 389
	}

	c, err := ldap.DialURL(fmt.Sprintf("%s://%s:%d", prot, cfg.Server, port), ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: !cfg.Verify}))
	if err != nil {
		return nil, err
	}

	err = c.Bind(cfg.Bind.Username, cfg.Bind.Password)
	if err != nil {
		return nil, err
	}

	if cfg.Verbose {
		fmt.Printf("# Connected to %s://%s:%d\n", prot, cfg.Server, port)
	}

	return c, nil
}

func GetGroups(conn *ldap.Conn, cfg Config) ([]*ldap.Entry, error) {
	var groups []*ldap.Entry

	for _, grp := range cfg.Groups {
		searchRequest := ldap.NewSearchRequest(
			cfg.BaseDN,
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0,
			cfg.Timeout,
			false,
			fmt.Sprintf("(memberOf=%s)", grp),
			cfg.Attributes,
			nil,
		)

		sr, err := conn.Search(searchRequest)
		if err != nil {
			continue
		}

		groups = append(groups, sr.Entries...)
	}

	if cfg.Verbose {
		fmt.Printf("# Found %d groups\n", len(groups))
	}

	return groups, nil
}

func SearchUser(conn *ldap.Conn, cfg Config, username string) (*User, error) {
	user := User{}

	for _, group := range cfg.Groups {
		userFilter := fmt.Sprintf("(%s=%s)", cfg.UserAttribute, username)
		groupFilter := fmt.Sprintf("(memberOf:1.2.840.113556.1.4.1941:=%s)", group)

		searchFilter := fmt.Sprintf("(&%s%s)", userFilter, groupFilter)

		searchRequest := ldap.NewSearchRequest(
			cfg.BaseDN,
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0,
			cfg.Timeout,
			false,
			searchFilter,
			cfg.Attributes,
			nil,
		)

		srUser, err := conn.Search(searchRequest)
		if err != nil {
			return nil, err
		}

		if len(srUser.Entries) == 0 {
			continue
		}

		if user.Login != "" {
			continue
		}

		user.Login = srUser.Entries[0].DN
		user.Username = srUser.Entries[0].GetAttributeValue("displayName")
		user.Group = IsAdmin(group, cfg)

		if (len(srUser.Entries)) > 1 {
			return nil, fmt.Errorf("too many entries returned")
		}
	}

	if cfg.Verbose {
		fmt.Printf("# Found user %s in group %s\n", user.Username, user.Group)
	}

	return &user, nil
}

func TryBind(conn *ldap.Conn, cfg Config, username string, password string) error {
	err := conn.Bind(username, password)
	if err != nil {
		if cfg.Verbose {
			fmt.Printf("# User %s failed to authenticate\n", username)
		}

		return err
	}

	if cfg.Verbose {
		fmt.Printf("# User %s authenticated\n", username)
	}

	conn.Unbind()
	return nil
}

func IsAdmin(ug string, cfg Config) string {
	isA := "system-users"

	if slices.Contains(cfg.Mappings["admin"], ug) {
		isA = "system-admin"
	}

	if cfg.Verbose {
		fmt.Printf("# User group %s is mapped to %s\n", ug, isA)
	}

	return isA
}
