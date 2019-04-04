# LDAP Auth Method

Gafaspot authenticates its users against an LDAP Server. Vault provides LDAP authentication through an [Auth Method](https://www.vaultproject.io/docs/auth/ldap.html). Gafaspot users are not meant to talk directly to the Vault server. But Gafaspot outsources the user authentication to Vault, which again performs LDAP authentication against some LDAP server. Therefore, you have to enable and configure Vault's LDAP Auth Method in the right way.

## Enable
You enable an Auth Method like this:

    curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @auth_ldap_enable.json http://127.0.0.1:8200/v1/sys/auth/ldap

Gafaspot expects the LDAP Auth Method to be enabled at path `auth/ldap`, so don't change the request path's end. To make the enabled Auth Method to be of type LDAP, following payload is needed:

    {
        "type": "ldap"
    }

## Configure
You upload a configuration with following command:
    
    curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @auth_ldap_config.json http://127.0.0.1:8200/v1/auth/ldap/config

An appropriate config would be something like:

    {
        "url": "ldaps://127.0.0.11:636",
        "userdn": "ou=Users,dc=example,dc=com",
        "groupdn": "ou=Groups,dc=example,dc=com",
        "groupfilter": "(&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))",
        "upndomain": "example.com"
    }

"url" should be of course your Active Directory's network address. If you want to connect via `ldaps`, make sure to upload the right certificate to the machine running Vault. "userdn" is the base DN under which to perform user search. "groupdn" is base to use for group membership search. With "userdn" and "groupdn" you locate the users which should be allowed to use Gafaspot. If you set "groupfilter" as in the example above, you enable LDAP to also resolute nested groups. "upndomain" defines a string which is appended to each user name in a login request. For example, a user's full login name as it is known by LDAP is usually something like userX@example.com, but the user will want to login only typing userX. So, put example.com into "upndomain".

## Map Policy
You will probably create an LDAP group for all users which should be allowed to use Gafaspot. Gafaspot needs to determine, whether authenticating users are members of this group. This is accomplished by configuring the LDAP Auth Method to assign a specific policy to members of this group. This policy's name is entered into Gafaspot's config file `gafaspot_config.yaml`. So, Gafaspot can check whether a authenticating user owns this policy.

Therefore, a new policy must be created:

    curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @policy_ldap_create.json http://127.0.0.1:8200/v1/sys/policy/gafaspot-user-ldap

where the last part of the request path is the policy's name inside Vault. Also write this name into `gafaspot_config.yaml`. The policy's content does not really matter. You can upload following payload to create an empty policy, only containing a comment:

    {
        "policy": "# This is an empty policy. Its i assigned to legitime Gafaspot users when authenticating with LDAP Auth Method. So, Gafaspot can recognize them by the policy name"
    }

Now, the policy needs to be mapped to the right LDAP group. This is done with following command:

    curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @auth_ldap_map_policy.json http://127.0.0.1:8200/v1/auth/ldap/groups/your_ldap_group_for_gafaspot_users

where the last part of the request path is the LDAP group's name, in which you want to put all Gafaspot users. The payload is the following:

    {
        "policies": "gafaspot-user-ldap"
    }



