# Explanations for Gafaspot Configuration
Besides setting up a Vault server, Gafaspot itself has to be configured.

All configuration for Gafaspot is read from one single config file: `gafapot_config.yaml`.
This file must be located in the same directory from which you run Gafaspot.
Create such a file by copying `example_config.yaml` which you find with the Gafaspot source code. Then adapt the file to reflect your desired settings.

Some config parameters have default values. This parameters are marked in the descriptions below. If present, this document uses the default as example value.

## Strutcure of config file
`gafapot_config.yaml` consists of three parts:

* general config for Gafaspot
* config concerning the database
* config concerning Vault

## General Config for Gafaspot
`webservice-address: 0.0.0.0:80` *(default value)*  
defines where the web server listens
___
`mailserver: mail.example.com:25`  
specifies the mail server (address and port) gafaspot can use to send e-mails to users. This feature is optional, so omit this configuration to disable emailing.
___

`gafaspot-mailaddress: gafaspot@gafaspot.com` *(default value)*  
defines a mail address under which gafaspot sends e-mails to its users. Gafaspot will not authenticate in any way, so the address does not have to exist. However, the mail server must allow sending unauthenticated mails.
___

`max-reservation-duration-days: 30` *(default value)*  
defines, how long one reservation for an environment is allowed to be (in days)
___
`max-queuing-time-months: 2`   *(default value)*  
defines, how far a reservation's start date can be in the future (in months)

## Config Concerning Database
`db-path: ./gafaspot.db`   *(default value)*  
specifies the file path of the SQLite database file Gafaspot will use. If database file does not yet exist, it will be created when starting Gafaspot.
___
`database-ttl-months: 12`   *(default value)*  
defines, how long a database entry is usually kept in the database after it is not used anymore. Currently, this ttl applies to the database tables 'users' and 'reservations'.  
For 'users', this means a user's table entry gets deleted if he has not logged in for this duration. Therefore, the deleted user will have to upload a new SSH public key if he wants to make reservations again.  
For the 'reservations' table the TTL specifies how long a reservation is kept after its expiry date before it will be deleted.
The value is given in months.  

## Config Concerning Vault
`vault-address: http://127.0.0.1:8200/v1`   *(default value)*  
network address of vault server. Beginning of each request path.  
Make sure to include the 'v1' ending which is currently the prefix for each route in vault. (reference: https://www.vaultproject.io/api/overview#http-api)
___
`approle-roleID: someID`  
`approle-secretID: someSecret`  
the credentials Gafaspot uses to authenticate against Vault. They are similar to a pair of username and password. You have to enable the approle auth method and create such credentials within Vault.
For more information, see the instructions about [Approle Auth Method](doc/auth_approle.md)  
___
`ldap-group-policy: gafaspot-user-ldap`   *(default value)*  
ldap-group-policy is the name of a Vault policy attached to tokens created with the LDAP Auth Method. When Gafaspot uses the LDAP Auth Method to verify its users, Vault requests over LDAP
* if the user credentials are valid at all and
* in case they are, to which groups the user belongs to.
Depending on the group, Vault associates preconfigured policies to the user and returns the policy names to Gafaspot. Based on this policy name Gafaspot decides whether the user is allowed to use Gafaspot or not.  
For more information about how to configure the LDAP Auth Method correctly, see the instructions about [LDAP Auth Method](doc/auth_ldap.md)
___
`environments:`  
The end of the Gafaspot config describes the composition of the different environments which you intend to manage with Gafaspot. Therefore, give a list of all environments at the first level like this:

```
    environments:

        demo0:
            ...

        demo1:
            ...

        demo2:
            ...

        ...
```

The environment's names are only allowed to contain **lowercase** ASCII letters, numbers and underscores. Don't use uppercase letters and blanks!  
Each environment has the following attributes: 

```
        demo0:
            show-name: DEMO 0
            description: "Some description for DEMO 0;
                          can use multiple lines and
                          HTML tags <br> for formatting."
            secrets-engines:
                ...
```

As you can see, you are able to provide an attribute `show-name` which is allowed to contain any character. This name will be displayed in web interface. Additionally, the web interface shows every instruction you write into `description`. Use HTML syntax for formatting. For example, you can include hyperlinks. You should explain in detail, which components are within the environment, which credentials to expect from the Secret Engines, and how the credentials map to the environments. `show-name` and `description` are optional.

Finally, you need to list all the Secrets Engines at the third level. Therefore, enable as many Secrets Engines in Vault as you need to perform credential changing for all devices in your environment. Additionally, enable one KV Secrets Engine for each credential-changing secrets engine. The Secrets Engines have to be enabled at the following paths:

    operate/<environment_name>/<secrets_engine_name>    => Some Secrets Engine offering new credentials
    store/<environment_name>/<secrets_engine_name>      => KV Secrets Engine which stores the credentials for the other Secrets Engine
In this example, environment_names would be `demo0`, `demo1` and so on. For enabling Secrets Engines at the right paths, read the [General Instructions about Secrets Engines](secengs_general.md).

Defining the Secrets Engines looks like this:

```
            ...
            secrets-engines:
                - name: NetApp
                  type: ontap
                  role: gafaspot
                
                - name: ActiveDirectory
                  type: ad
                  role: gafaspot

                - ...
```
                
The Secrets Engine's name may only contain ASCII letters, numbers and underscores. Anyway, try to choose a descriptive name, as this name will be displayed in web interface when user request credentials. As described in [Secrets Engines General](secengs_general.md), the name is the last part of the path, under which you enable the Secrets Engine in Gafaspot.  
`type` is one of:
* ad
* ssh
* database
* ontap

You do not have to explicitly mention KV Secrets Engines in the config file, as they are always related to another Secrets Engine.

`role` is the name of the role you configure with the respective Secrets Engine. How you create the role is described in the respective instructions about the Secrets Engine type.

---
*Go to [next page](database_scheme.md)...*  
*Go to [table of contents](README.md)...*
