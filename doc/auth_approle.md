# AppRole Auth Method

Gafaspot needs to access Vault for performing operations like changing credentials. Therefore, Gafaspot uses the AppRole Auth Method to retrieve valid vault tokens when needed. So you have to set this Auth Method up within Vault.

To log in with AppRole Auth Method, Gafaspot needs to present credentials to Vault, similar to username and password. You will generate such credentials and copy them to gafaspot_config.yaml. This of course somewhat compromises the extra security provided by Vault's storage encryption -- in the end, Gafaspot's secrets are again only protected by the operating system on which you run Gafaspot. But as Gafaspot is not meant to be a high security system, this is the simplest approach.

At least, a better security can be accomplished by slightly extending Gafaspot. For example, you could make it to rotate its credentials at start, so they are not as plain text in persistent storage anymore. Most of such security improvements would need some kind of interaction with Gafaspot or Vault when restarting Gafaspot. As this is kind of disruptive, such a mechanism is not implemented yet. But you can still delete the AppRole credentials from gafaspot_config.yaml after you started Gafaspot, as the file is only read once at program start.

## Enable
You enable an Auth Method like this:

    curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @auth_approle_enable.json http://127.0.0.1:8200/v1/sys/auth/approle

Gafaspot expects the AppRole Auth Method to be enabled at path `auth/approle`, so don't change the request path's end. To make the enabled Auth Method to be of type AppRole, following payload is needed:

    {
        "type": "approle"
    }

## Create Policy
To give Gafaspot the proper permissions inside Vault, a policy is needed. An appropriate policy can be found in file [policy_gafaspot_approle.hcl](policy_gafaspot_approle.hcl). Upload it with following command:

    curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @policy_approle_create.json http://127.0.0.1:8200/v1/sys/policy/gafaspot-approle

The last part of the request path is the policy's name. The payload policy_approle_create.json is not copied to this page because of its length, but it is included in the json_payload directory.

## Create Role
There is no write to path `/config` for this Auth Method. Instead you have to create a role:

    curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @auth_approle_role.json http://127.0.0.1:8200/v1/auth/approle/role/gafaspot

The last part of the request path is the role name. Use following payload:

    {
        "policies": ["gafaspot-approle"],
        "secret_id_num_uses": 0,
        "secret_id_ttl": "",
        "token_num_uses": 0,
        "token_ttl": "1m"
    }

With parameter `policies`, you map the just created policy to tokens retrieved with this role. `secrets_id_num_uses` and `secret_id_ttl` restrict, how often and how long the same credentials can be used. The value `0` respective `""` means the use is unlimited, which is what we want for Gafaspot. `token_num_uses` and `token_ttl` determine the same for each token, retrieved by applying the Auth Method one time. Gafaspot will fetch a token before each series of actions in Vault, so restrict the ttl to a small amount of time. One minute should fit well.

## Retrieve fresh credentials
Credentials for the AppRole Auth Method are pairs of two strings called `role-id` and `secret-id`. They work in the same way as a pair of username and password. You need to reed both values from Vault.

Read the role's `role-id` like this:

    curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' http://127.0.0.1:8200/v1/auth/approle/role/gafaspot/role-id

Then, generate a `secret-id`:

    curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST http://127.0.0.1:8200/v1/auth/approle/role/gafaspot/secret-id

It has to be a POST request, but its payload can be empty.

Copy `role-id` and `secret-id` from the request's answers to gafaspot_config.yaml

---
*Go back to [table of contents](README.md)...*