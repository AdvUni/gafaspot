# Ontap Secrets Engine

The Ontap Secrets Engine is a custom plugin which is not available in the Vault repository. It can change passwords for a NetApp device running with the operating system ONTAP. It changes ONTAP passwords using the XML-RPC API ONTAPI.

As this Secrets Engine was developed specifically for the use with Gafaspot, its functionality is limited to what is needed by Gafaspot.

## Enable
After the ontap plugin is [registered](https://www.vaultproject.io/docs/internals/plugins.html#plugin-registration) within Vault, you can enable the Secrets Engine like this:

```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @ontap_enable.json http://127.0.0.1:8200/v1/sys/mounts/operate/<environment_name>/NetApp
```

with the payload:

```json
{
    "type": "ontap"
}
```

Also enable a respective KV storage Secrets Engine:

```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @kv_enable.json http://127.0.0.1:8200/v1/sys/mounts/store/<environment_name>/NetApp
```

which has the adapted payload:

```json
{
    "type": "kv",
    "version": 1
}
```

## Configure
You can upload a configuration with the following command:
    
```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @ontap_config.json http://127.0.0.1:8200/v1/operate/<environment_name>/NetApp/config
```

An appropriate config would be something like:

```json
{
    "url":"https://127.0.0.1/servlets/netapp.servlets.admin.XMLrequest_filer",
    "vaultUser": "admin_vault",
    "vaultUserPass": "Password123"
}
```

"url" should be your NetApp's network address. Extend it to the path on which ONTAPI listens. "vaultUser" and "vaultUserPass" have to be filled with username and password of a local ONTAP user account which has enough permissions to change passwords for other user accounts. The Secrets Engine will use these credentials to authenticate against ONTAP.

## Create Role
With a role you tell the Secrets Engine *for which* account it should change passwords. Create a role with the following command:

```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @ontap_role.json http://127.0.0.1:8200/v1/operate/<environment_name>/NetApp/roles/gafaspot
```

The last part of the url is the role name, used by the Secrets Engine to identify the ONTAP account. You can choose it freely, but you have to specify it in the gafaspot_config.yaml file as well.
The command's payload looks like this:

```json
{
    "ontap_name": "gafaspot_on_ontap"
}
```

Change ontap_name to the user name of any local ONTAP account you want to manage with Gafaspot.

---
*Go to [next page](config_explanation.md)...*  
*Go to [table of contents](README.md)...*
