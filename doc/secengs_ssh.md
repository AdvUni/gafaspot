> Note: There is now an alternative SSH Secrets Engine implementation, which not needs the certificates anymore. The original SSH Secrets Engine is still supported, but the [SSH-Pubkey](secengs_sshpubkey.md) might fit your needs better.

# SSH Secrets Engine
The [SSH Secrets Engine](https://www.vaultproject.io/docs/secrets/ssh/index.html) is available in three different modes. Gafaspot only supports the [Signed Certificates Mode](https://www.vaultproject.io/docs/secrets/ssh/signed-ssh-certificates.html). In this mode, the Secrets Engine does not interact directly with any external device which makes configuration relatively easy. On the other hand, you will have to do some configuration on the target machines to make them accept the Secrets Engine's certificates. 

The Secrets Engine is applicable to manage access to all accounts allowing ssh login using authentication with a signed private key. This will apply primarily to unix-like operating systems. If you have several of such devices in one environment, you can configure them all to accept certificates of the same instance, so you will need only one SSH Secrets Engine per environment.

The signing process works like this: When reserving a environment with an SSH Secrets Engine, Gafaspot performs a `sign` request towards the Secrets Engine, containing the user's public ssh key. Therefore, the user must have uploaded such a key in the Gafaspot web interface previously. The Secrets Engine has its own ssh key pair and signs the user key, meaning it returns a certificate confirming the user key for a specified validity period. The user can retrieve the certificate in the web interface of Gafaspot. Now he can login with his own private key and the certificate to every SSH machine which is configured to trust the Secrets Engines certificates.

## Enable
You can enable the Secrets Engine like this:

```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @ssh_enable.json http://127.0.0.1:8200/v1/sys/mounts/operate/<environment_name>/SSH
```

with this payload:

```json
{
    "type": "ssh"
}
```

Also enable a respective KV storage Secrets Engine:

```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @kv_enable.json http://127.0.0.1:8200/v1/sys/mounts/store/<environment_name>/SSH
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
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @ssh_config.json http://127.0.0.1:8200/v1/operate/<environment_name>/SSH/config/ca
```

As you can see, the last part of the request path is `ca`. This defines that the SSH Secrets Engine will be used with the Signed Certificates Mode. Furthermore, use the following config:

```json
{
    "generate_signing_key": true
}
```

The Secrets Engine will generate a new ssh key pair and return the public key to you. It is your job to register it in all your environment's devices as Certificate Authority. Later, you can retrieve the key again with:

```sh
curl http://127.0.0.1:8200/v1/operate/<environment_name>/SSH/public_key
```

## Create Role
The concept of roles is a bit strange with the signed certificates mode, but you need one just the same. The command is:

```sh
curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @ssh_role.json http://127.0.0.1:8200/v1/operate/<environment_name>/SSH/roles/gafaspot
```

The last part of the url is the role name which you also need to specify in the config file gafaspot_config.yaml.
The following payload should work:

```json
{
    "key_type": "ca",
    "allow_user_certificates": true,
    "default_extension": [
        {
            "permit_pty": ""
        }
    ]
}
```
The content of the `default_extension` field is necessary to show a prompt after logging in with the signed ssh certificates (See [Vault Documentation](https://www.vaultproject.io/docs/secrets/ssh/signed-ssh-certificates.html#no-prompt-after-login)).

## Register SSH Secrets Engine as Certificate Authority
This step depends on the target operating system for which you want to manage accounts with Gafaspot.

On unix-like operating systems, there is usually a file called `authorized_keys` somewhere. This file holds public ssh keys, whose owners are allowed to log in to the system. In here, one can also register keys which are not meant to directly log in, but to sign other keys to give them the permission. This is what we want for the SSH Secrets Engine.

As already stated, you can read the Secrets Engine's public key like this:

```sh
curl http://127.0.0.1:8200/v1/operate/<environment_name>/SSH/public_key
```

Find the `authorized_keys` file, start a new line and write following statements, separated by blanks:

```
cert-authority ssh-rsa AAAAB3NzaC...
```

where the third part is the Secrets Engine's key of course.

If this configuration won't allow you to log in with signed keys out of the box, you may need to adapt the file `/etc/ssh/sshd_config`, which contains general settings about ssh login. Reload the file with

```sh
sudo service sshd restart
```

after you made changes.


---
*Go to [next page](secengs_sshpubkey.md)...*  
*Go to [table of contents](README.md)...*
