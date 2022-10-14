# LDAP Auth Method

Gafaspot authenticates its users against an LDAP Server. Vault provides LDAP authentication through an [Auth Method](https://www.vaultproject.io/docs/auth/ldap.html). Gafaspot users are not meant to talk directly to the Vault server. However, Gafaspot outsources the user authentication to Vault, which again performs LDAP authentication against some LDAP server. Therefore, you have to enable and configure Vault's LDAP Auth Method correctly.

## Enable
You can enable an Auth Method like this:

```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @auth_ldap_enable.json http://127.0.0.1:8200/v1/sys/auth/ldap
```

Gafaspot expects the LDAP Auth Method to be enabled at path `auth/ldap`, so don't change the end of the request path. To configure the enabled Auth Method to be of type LDAP, following payload is needed:

```json
{
    "type": "ldap"
}
```

## Configure
You can upload a configuration with the following command:

```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @auth_ldap_config.json http://127.0.0.1:8200/v1/auth/ldap/config
```

An appropriate config would be something like:

```json
{
    "url": "ldaps://127.0.0.11:636",
    "userdn": "ou=Users,dc=example,dc=com",
    "groupdn": "ou=Groups,dc=example,dc=com",
    "groupfilter": "(&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))",
    "upndomain": "example.com"
}
```

"url" should be your LDAP or Active Directory Domain Controller's network address. If you want to connect via `ldaps` (using TLS), make sure to upload the right server certificate to the machine running Vault. "userdn" is the base DN under which to perform user search. "groupdn" is the base DN to use for group membership search. With "userdn" and "groupdn" you locate the users which should be allowed to use Gafaspot. If you set a "groupfilter", as in the example above, you enable LDAP to also resolve nested groups. "upndomain" defines a string which is appended to each user name in a login request. For example, a user's full login name as it is known by LDAP is usually something like userX@example.com, but the user will want to login only typing userX. In this case you would put `example.com` into "upndomain".

## Map Policy
You will probably want to create an LDAP group for all users which should be allowed to use Gafaspot. Gafaspot needs to determine whether authenticated users are members of this group. This is accomplished by configuring Vault's LDAP Auth Method to assign a specific policy to members of this group. This policy's name is entered into Gafaspot's config file `gafaspot_config.yaml`. So, Gafaspot can check whether a authenticating user owns this policy.

Therefore, a new policy must be created:

```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @policy_ldap_create.json http://127.0.0.1:8200/v1/sys/policy/gafaspot-user-ldap
```

Here, the last part of the request path (`gafaspot-user-ldap`) is the policy's name inside Vault. You will also need to write this name into `gafaspot_config.yaml`. The policy's content does not really matter. You can upload the following payload to create an empty policy only containing a comment:

```json
{
    "policy": "# This is an empty policy. It is assigned to legitimate Gafaspot users when authenticating with the LDAP Auth Method so that Gafaspot can recognize them by the policy name"
}
```

Now, the policy has to be mapped to the right LDAP group. This is done with the following command:

```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @auth_ldap_map_policy.json http://127.0.0.1:8200/v1/auth/ldap/groups/your_ldap_group_for_gafaspot_users
```

where the last part of the request path is the LDAP group's name in which you want to put all Gafaspot users. The payload is the following:

```json
{
    "policies": "gafaspot-user-ldap"
}
```


---
*Go to [next page](secengs_general.md)...*  
*Go to [table of contents](README.md)...*
