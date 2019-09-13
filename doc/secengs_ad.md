# AD Secrets Engine

The [Active Directory Secrets Engine](https://www.vaultproject.io/docs/secrets/ad/index.html) can change passwords for accounts in an Active Directory. Therefore, you have to enable it with the correct path, configure it to communicate with the right AD and create a role which specifies for which accounts passwords can be changed.

You can configure a TTL which specifies how long passwords set by this Secrets Engine shall be valid. But the way the AD Secrets Engines treats those TTLs is a bit strange. No Leases are created, instead the Secrets Engine performs something called 'Lazy Rotation'. Basically, this means that the Secrets Engine doesn't 'automatically' change passwords at all. It only checks the TTL if you do a `creds` request. If the TTL is expired, the Secrets Engine creates a new password, otherwise it sends back the old.

This behavior does not really fit to the needs of Gafaspot. Either, the passwords should expire automatically in the background when a reservation is over, or Gafaspot needs a reliable method of triggering a password change explicitly. To make the Secrets Engine behave in a usable way, Gafaspot needs you to set the Secrets Engine's default TTL to a very small value, which will cause every `creds` request from Gafaspot to change the password.


## Enable
You enable the Secrets Engine like this:

```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @ad_enable.json http://127.0.0.1:8200/v1/sys/mounts/operate/<environment_name>/ActiveDirectory
```

with payload:

```json
{
    "type": "ad"
}
```

Also, enable a respective KV storage Secrets Engine:

```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @kv_enable.json http://127.0.0.1:8200/v1/sys/mounts/store/<environment_name>/ActiveDirectory
```

which has the adapted payload:

```json
{
    "type": "kv",
    "version": 1
}
```


## Configure
You can upload a configuration with following command:

```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @ad_config.json http://127.0.0.1:8200/v1/operate/<environment_name>/ActiveDirectory/config
```

An appropriate config would be something like:

```json
{
    "url": "ldaps://127.0.0.11:636",
    "binddn": "cn=Administrator,cn=Users,dc=example,dc=com",
    "bindpass": "Password123",
    "userdn": "ou=Users,dc=example,dc=com",
    "ttl": "1",
    "max_ttl": "1"
}
```

"url" should be your Active Directory's network address. You will have to use the `ldaps` protocol instead of `ldap` since the Active Directory denies credential changing over insecure connections. Make sure to upload the Active Directories certificate to the machine running Vault. Otherwise, you can use the `insecure_tls` flag within the configuration. See the [API documentation](https://www.vaultproject.io/api/secret/ad/index.html) for more information.
For "binddn" and "bindpass" you have to fill in username and password of an AD user which can be used by Vault to perform password changes for other users. Make sure that the user has all required permissions inside the AD. "userdn" is the base DN under which to perform user search -- here you can determine where the accounts are found for which the password changes should be performed later on. The TTL parameters must be set to a very short time as discussed above.

Be aware, that the AD Secrets Engine will not give you any feedback about whether the configuration is valid after uploading it.

## Create Role
With a role you tell the Secrets Engine *for which* account it should change passwords. Create a role with the following command:

```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @ad_role.json http://127.0.0.1:8200/v1/operate/<environment_name>/ActiveDirectory/roles/gafaspot
```

The last part of the url is the role name. This is the name used by the Secrets Engine to identify the AD account. You can choose it freely, but you have to specifiy it in the gafaspot_config.yaml file as well.
The command's payload looks like this:

```json
{
    "service_account_name": "gafaspot_on_ad@example.com"
}
```

Change the service_account_name to the name of an AD user which you want to manage with Gafaspot. It does not have to be something like a service account. Normal user accounts work as well.

## Test Setup
As Vault might not have told you if some of the configuration is invalid, you should perform a `creds` request to test whether your setup works as expected:

```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' http://127.0.0.1:8200/v1/operate/<environment_name>/ActiveDirectory/creds/gafaspot
```

---
*Go to [next page](secengs_ssh.md)...*  
*Go to [table of contents](README.md)...*
