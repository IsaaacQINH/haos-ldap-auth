package gateway

import (
	"fmt"
	"slices"

	"github.com/go-ldap/ldap/v3"
)

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

	c, err := ldap.DialURL(fmt.Sprintf("%s://%s:%d", prot, cfg.Server, port))
	if err != nil {
		return nil, err
	}

	err = c.Bind(cfg.Bind.Username, cfg.Bind.Password)
	if err != nil {
		return nil, err
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

	return groups, nil
}

func SearchUser(conn *ldap.Conn, cfg Config, groups []string, username string) (*ldap.Entry, error) {
	userFilter := fmt.Sprintf("(%s=%s)", cfg.UserAttribute, username)
	groupFilter := "(|"

	for _, grp := range groups {
		groupFilter += fmt.Sprintf("(memberOf=%s)", grp)
	}

	groupFilter += ")"
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

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	if (len(sr.Entries)) > 1 {
		return nil, fmt.Errorf("too many entries returned")
	}

	return sr.Entries[0], nil
}

func TryBind(conn *ldap.Conn, cfg Config, username string, password string) error {
	err := conn.Bind(username, password)
	if err != nil {
		return err
	}

	conn.Unbind()
	return nil
}

func IsAdmin(entry *ldap.Entry, cfg Config) string {
	isA := "system-users"

	entryGroups := entry.GetAttributeValues("memberOf")
	cts := slices.Contains(entryGroups, cfg.Mappings["admin"])

	if cts {
		isA = "system-admin"
	}

	return isA
}
