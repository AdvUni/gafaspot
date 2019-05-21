# Database Secrets Engine

The [Database Secrets Engine](https://www.vaultproject.io/docs/secrets/databases/index.html) unites access management for several databases. The Vault documentation lists the suppoted database. In Theory, one Database Secrets Engine can handle connections to a number of different databases. So if you have multiple databases in one environment, you can probably handle them all with the same Secrets Engine. But to create a less complicated setup it is better to enable one Secrets Engine per database.

When performing a `creds` request against a Database Secrets Engine it creates a [Lease](https://www.vaultproject.io/docs/concepts/lease.html) containing a [Dynamic Secret](https://www.hashicorp.com/blog/why-we-need-dynamic-secrets). This means, that the retrieved credentials are only valid for a specific period of time and the Secrets Engine revokes them automatically in the background when expired (in contrast to the 'Lazy Rotation' concept which is implemented with the Active Directory Secrets Engine). Unfortunately, the Database Secrets Engine offers no way of defining the duration on a per credential basis. Instead it is determined within the Secret Engine's or the role's configuration. As Gafaspot needs to declare the validity period individually for each reservation Gafaspot kind of bypasses the logic of Leases. Gafaspot sets the default lease duration to the maximal reservation duration given in gafaspot_config.yaml. When starting a reservation, Gafaspot requests a new lease from the Database Secrets Engine. At a reservation's ending, Gafaspot manually revokes the lease. Therefore, it is important you do not define the default lease duration yourself when configuring the Secrets Engine or a role for it. Otherwise, it might not be possible to make long-term reservations.

## Enable
Enable the Secrets Engine like this:

```
    curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @database_enable.json http://127.0.0.1:8200/v1/sys/mounts/operate/<environment_name>/DB
```

with payload:

```
    {
        "type": "database"
    }
```

Also enable a respective KV storage Secrets Engine:

```
    curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @kv_enable.json http://127.0.0.1:8200/v1/sys/mounts/store/<environment_name>/DB
```

which has the adapted payload:

```
    {
        "type": "kv",
        "version": 1
    }
```

## Configure
With uploading the config you determine which database type you serve. As one Database Secrets Engine can handle connections to multiple databases at once, you have to establish a name for your specific database configuration, which is given as the last parameter in the request url.

This guide limits itself to describe the configuration with a [MYSQL database](https://www.vaultproject.io/docs/secrets/databases/mysql-maria.html). For other databases, see the Vault documentation.

You can upload a configuration with the following command:
    
```
    curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @database_config.json http://127.0.0.1:8200/v1/operate/<environment_name>/DB/config/my_database
```

For MYSQL, the config would be something like:

```
{
    "plugin_name": 		"mysql-database-plugin",
    "allowed_roles": 	"*",
    "connection_url": 	"{{username}}:{{password}}@tcp(127.0.0.1:3306)/",
    "username": 		"admin_vault",
    "password": 		"Password123"
}
```

"plugin_name" defines the database you want to handle. You will probably create only one role anyway, so you can set "allowed_roles" to all. The "connection_url" does not only contain the database's network address, but the whole Data Source Name. You can probably copy it as it is. "username" and "password" are the credentials of an existing database user which has enough permissions to create and remove other users. The Secrets Engine will use these credentials to authenticate against the database.

## Create Role
The Database Secrets Engine does not change the password for an existing account if requested. Instead, when performing a `creds` request, it creates a new user, which is removed again after some time. A role inside the Database Secrets Engine defines which properties such a newly created user will have.
Create a role with following command:

```
    curl --header 'X-Vault-Token: '"$VAULT_TOKEN"'' --request POST --data @database_role.json http://127.0.0.1:8200/v1/operate/<environment_name>/DB/roles/gafaspot
```

The last part of the url is the role name which you also need to specify in the config file gafaspot_config.yaml.
The following payload should work:

```
    {
        "db_name": "mysql",
        "creation_statements": ["CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}'", "GRANT ALL ON *.* TO '{{name}}'@'%' WITH GRANT OPTION"]
    }
```

"db_name" defines again for which kind of database this role is created. "creation_statements" is a list of statements which the Secrets Engines executes when creating the new user. This needs to be set explicitly, because this is the only point where it is possible to determine which permissions the new user will have within the database. The statement `GRANT ALL ON *.* TO '{{name}}'@'%' WITH GRANT OPTION` should give all permissions to users created with the Secrets Engine. It is also possible to define "revocation_statements", but this is not required. It defaults to just deleting the user without any further actions.

---
*Go to [next page](secengs_ontap.md)...*  
*Go to [table of contents](README.md)...*
