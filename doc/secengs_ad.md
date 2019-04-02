# AD Secrets Engine

The [Active Directory Secrets Engine](https://www.vaultproject.io/docs/secrets/ad/index.html) can change credentials for accounts in an Active Directory.

You enable it like this:

    curl --header "X-Vault-Token: ..." --request POST --data @ad_enable.json http://127.0.0.1:8200/v1/sys/mounts/operate/<environment_name>/ActiveDirectory

with payload:

    {
        "type": "ad"
    }

Don't forget to enable a respective KV storage Secrets Engine:

    curl --header "X-Vault-Token: ..." --request POST --data @kv_enable.json http://127.0.0.1:8200/v1/sys/mounts/operate/<environment_name>/ActiveDirectory

which has the adapted payload:

    {
        "type": "kv",
        "version": 1
    }
