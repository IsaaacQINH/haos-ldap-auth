# 🟢 haos-ldap-auth

A tiny Go-based executable to bridge **LDAP Authentication** into your **Home Assistant OS** setup — for anyone who wants their smart home users managed centrally and cleanly via LDAP.

## 🌟 What is this?

This tool lets you use your LDAP server to authenticate users logging in to Home Assistant OS. It's lightweight, simple, and avoids bloating your HA instance.

## 🚀 Quickstart

1. **Clone this repo**

```bash
git clone https://github.com/IsaaacQINH/haos-ldap-auth
cd haos-ldap-auth
```

2. **Edit your config**
```yml
server: ldap.example.com
tls: false
basedn: DC=base,DC=example,DC=com
bind:
  user: LDAP_BIND_USER_DN
  password: LDAP_BIND_USER_PASSWORD
groups:
  - CN=GROUP1,OU=Groups,DC=base,DC=example,DC=com
  - CN=GROUP2,OU=Groups,DC=base,DC=example,DC=com
mappings:
  admin: 
    - CN=Role,OU=Roles,DC=base,DC=example,DC=com
user_attribute: sAMAccountName
attributes: [sAMAccountName, displayName, memberOf]
timeout: 10
```

3. **Build and run**

```bash
go build -o haos-ldap-auth .
./haos-ldap-auth
```

4. **Add this to your Home Assistant `configuration.yml`**

Note: Ensure `meta` is set true - otherwise you won't able to set a username/group for new created users

```yml
homeassistant:
    auth_providers:
        - type: command_line
          name: 'LDAP Auth'
          command: '/path/to/haos-ldap-auth'
          meta: true
        - type: homeassistant
```

Done!

