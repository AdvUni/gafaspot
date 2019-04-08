# AppRole Auth Method

Gafaspot needs to access Vault for performing operations like changing credentials. Therefore, Gafaspot uses the AppRole Auth Method to retrieve valid vault tokens when needed. So you have to set this Auth Method up within Vault.

To log in with AppRole Auth Method, Gafaspot needs to present credentials to Vault, similar to username and password. You will generate such credentials and copy them to gafaspot_config.yaml. This of course somewhat compromises the extra security provided by Vault's storage encryption -- in the end, Gafaspot's secrets are again only protected by the operating system on which you run Gafaspot. But as Gafaspot is not meant to be a high security system, this is the simplest approach.

At least, a better security can be accomplished by slightly extending Gafaspot. For example, you could make it to rotate its credentials at start, so they are not as plain text in persistent storage anymore. Most of such security improvements would need some kind of interaction with Gafaspot or Vault when restarting Gafaspot. As this is kind of disruptive, such a mechanism is not implemented yet.

## Enable
You enable an Auth Method like this:

    curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @auth_approle_enable.json http://127.0.0.1:8200/v1/sys/auth/approle

Gafaspot expects the AppRole Auth Method to be enabled at path `auth/approle`, so don't change the request path's end. To make the enabled Auth Method to be of type AppRole, following payload is needed:

    {
        "type": "approle"
    }

## Create Policy
To give Gafaspot the proper permissions inside Vault, a policy is needed. An appropriate policy can be found in file `policy_gafaspot_approle.hcl`. Upload it with following command:

    curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @policy_approle_create.json http://127.0.0.1:8200/v1/sys/policy/gafaspot-approle

and following payload:

    

## Create Role
There is no write to path `/config` for this Auth Method. Instead you have to create a role:

    curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @auth_approle_role.json http://127.0.0.1:8200/v1/operate/auth/approle/roles/gafaspot

