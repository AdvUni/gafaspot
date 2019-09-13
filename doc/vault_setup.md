# Vault Setup

Read through Vault's [Getting Started Guide](https://learn.hashicorp.com/vault/). Follow the instructions to become familiar with Vault. Start your own Vault Server. Create a config.hcl file which fits your needs. When choosing a storage backend for the Vault Server, type [file](https://www.vaultproject.io/docs/configuration/storage/filesystem.html) seems to be the simplest one and should probably work fine with Gafaspot. Furthermore, the following guide assumes that you run Vault and Gafaspot on the same machine and therefore, Vault's API listens on localhost on port 8200.

After starting the server, initialize it with the following command:

```sh
curl --request PUT --data @vault_init.json http://127.0.0.1:8200/v1/sys/init
```

The file `vault_init.json` can be found together with many other JSON payload snippets in the sub directory [`json_payload/`](json_payload).
The contents of `vault_init.json` are:

```json
{
    "secret_shares": 1,
    "secret_threshold": 1
}
```

This will return you a single unseal key together with a root token. Gafaspot is not meant to manage super sensible secrets, so there is probably no need to split the responsibility for the unsealing process to several persons. You will need the unseal key to unlock Vault each time you restart it. The root token is necessary to supply yourself with access rights to Vault, which are needed to perform any configuration. So note both values down and keep them somewhere.

You can save the root token to the environment variable `$VAULT_TOKEN` to simplify the copying of further commands from this guide:

```sh
export VAULT_TOKEN='s.3eX...'
```

Unseal Vault with:

```sh
curl --request PUT --data @vault_unseal.json http://127.0.0.1:8200/v1/sys/unseal
```

`vault_unseal.json` is not included in the collection of JSON payload snippets in this repository as it contains your unseal key which is not known by this guide. It looks like:

```json
{
    "key": "abcd1234..."
}
```

From now on, Vault is up and ready for interaction. To be used with Gafaspot, you need to enable and configure some Auth Methods and many Secrets Engines. See the respective pages for more details:

* [AppRole Auth Method](auth_approle.md)
* [LDAP Auth Method](auth_ldap.md)
* [Secrets Engines](secengs_general.md)
    * [Active Directory Secrets Engine](secengs_ad.md)
    * [SSH Secrets Engine (Signed Certificates)](secengs_ssh.md)
    * [Database Secrets Engine](secengs_database.md)
    * [Ontap Secrets Engine](secengs_ontap.md)

---
*Go to [next page](auth_approle.md)...*  
*Go to [table of contents](README.md)...*
