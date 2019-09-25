# SSH-Pubkey Secrets Engine

The SSH-Pubkey Secrets Engine is a custom plugin which is available at [GitHub](github.com/AdvUni/vault-plugin-secrets-ssh-pubkey). It is a variation from Vault's builtin [SSH Secrets Engine](https://www.vaultproject.io/docs/secrets/ssh/index.html). The SSH-Pubkey plugin is kind of a combination of the [Dynamic Keys Mode](https://www.vaultproject.io/docs/secrets/ssh/dynamic-ssh-keys.html) and the [Signed Certificates Mode](https://www.vaultproject.io/docs/secrets/ssh/signed-ssh-certificates.html): It physically writes ssh public keys into a target machine's authorized_keys file like in the dynamic mode, but instead of generating those keys itself first, it receives the key by the user as in the certificate mode. As a result, one can simply use its own ssh key pair to log into the target machine after using the Secrets Engine.

Like the SSH Secrets Engine, the SSH-Pubkey Secrets Engine handles unix-like accounts which allows ssh logins via public key authentication. The Secrets Engine is designed to work with one target engine, so you need to enable one instance for each SSH account in your environment.

When performing a `creds` request against the SSH-Pubkey Secrets Engine it creates a [Lease](https://www.vaultproject.io/docs/concepts/lease.html) containing a [Dynamic Secret](https://www.hashicorp.com/blog/why-we-need-dynamic-secrets). This means, that the retrieved access is only valid for a specific period of time and the Secrets Engine revokes it automatically in the background when expired (in contrast to the 'Lazy Rotation' concept which is implemented with the Active Directory Secrets Engine). Vault also revokes a lease if the token with which it was created expires. So, at reservation start, Gafaspot creates an orphan vault token with a life span matching the reservation duration. At a reservation's ending, this token expires and Vault revokes all leases associated with it.

To not get Leases revoked to early, Gafaspot sets the default lease duration to the maximal reservation duration given in gafaspot_config.yaml at startup. Hence it is important you do not define the default lease duration yourself when configuring the Secrets Engine or a role for it. Otherwise, it might not be possible to make long-term reservations.

## Enable
As soon the SSH-Pubkey plugin is [registered](https://www.vaultproject.io/docs/internals/plugins.html#plugin-registration) within Vault, you can enable the Secrets Engine like this:

```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @sshpubkey_enable.json http://127.0.0.1:8200/v1/sys/mounts/operate/<environment_name>/SSH-Pubkey
```

with the payload:

```json
{
    "type": "ssh-pubkey"
}
```
This assumes that you registered the plugin with the name "ssh-pubkey" in vault.

Also enable a respective KV storage Secrets Engine:

```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @kv_enable.json http://127.0.0.1:8200/v1/sys/mounts/store/<environment_name>/SSH-Pubkey
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
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @sshpubkey_config.json http://127.0.0.1:8200/v1/operate/<environment_name>/SSH-Pubk/config
```

An appropriate config would be something like:

```json
{
    "url": "127.0.0.1",
    "private_key": "-----BEGIN RSA PRIVATE KEY-----\n ...",
    "public_key": "ssh-rsa AAAAB3NzaC1yc2EAA..."
}
```

The configuration defines, how Vault should connect to the target machine.
"url" should be your target machine's network address. It can either be an URL or an IP address.
Additionally, you have to pass a key pair which has privileged access at the target machine, which means the private key is installed in the file `/root/.ssh/authorized_keys`. The Secrets Engine will authenticate with this key to add and remove keys in the `authorized_keys` files of other system users.

Note, that private keys usually contain line breaks which are not allowed in json strings. So, you first must encode them with \n.

The secrets Engine access the `authorized_keys` files with a shell script. If your target machine is not a classic linux host, you might need individual commands instead. Therefore, you have the ability to provide the `install_script` parameter during configuration. For reference, see the [default script](https://github.com/AdvUni/vault-plugin-secrets-ssh-pubkey/blob/master/plugin/linux_install_script.go).

## Create Role
With a role you tell the Secrets Engine *for which* system user it should modify the authorized keys. Create a role with the following command:

```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @sshpubkey_role.json http://127.0.0.1:8200/v1/operate/<environment_name>/SSH-Pubkey/roles/gafaspot
```

The last part of the url is the role name, used by the Secrets Engine to identify the system user. You can choose it freely, but you have to specify it in the gafaspot_config.yaml file as well.
The command's payload looks like this:

```json
{
    "username": "example-user"
}
```

Change `example-user` to the target machine's user which you want to manage with Gafaspot.

---
*Go to [next page](secengs_database.md)...*  
*Go to [table of contents](README.md)...*