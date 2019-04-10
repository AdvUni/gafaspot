# Explanations for Gafaspot Configuration
All configuration for Gafaspot is read from one single config file: `gafapot_config.yaml`.
This file must be located in the same directory from which you run Gafaspot.
Create such a file by copying `example_config.yaml` which you find withing Gafaspot source code. Then adapt the file to your setting.

## Strutcure of config file
`gafapot_config.yaml` consists of three parts:

* general config for Gafaspot
* config concerning database
* config concerning Vault

## General Config for Gafaspot
`webservice-address: 127.0.0.1:8080`  
defines where the web server listens
___
`max-reservation-duration-days: 30`  
defines, how long one reservation for an environment is allowed to be (in days)
___
`max-queuing-time-months: 2`  
defines, how far a reservation's start can be in the future (in months)

## Config Concerning Database
`db-path: ./gafaspot.db`  
specifies the file path of the SQLite database file Gafaspot will use. If database file doesn't yet exist, it will be created with starting Gafaspot.
___
`database-ttl-months: 12`  
defines, how long a database entry is usually kept in database, after it is not used anymore. Currently, this ttl applies to the database tables 'users' and 'reservations'.  
For 'users', this means a user's table entry gets deleted if he haven't logged in for this duration. The only affect of this is, that the deleted user will have to upload a new SSH public key, if he wants to do reservations again.  
For table reservations, the TTL means a reservation gets deleted when this time has elapsed after the reservation's expiring date.  
Value is given in months.  

## Config Concerning Vault
`vault-address: http://127.0.0.1:8200/v1`  
network address of vault server. Beginning of each request path.  
Make sure to include the 'v1' ending which is currently the prefix for each route in vault. (reference: https://www.vaultproject.io/api/overview)
___
`approle-roleID: someID`  
`approle-secretID: someSecret`  
the credentials, Gafaspot uses to authenticate against Vault. They are similar to a pair of username and password. You have to enable the approle auth method and create such credentials inside Vault.
For more information, see the instructions about [Approle Auth Method](doc/auth_approle.md)  
___
`ldap-group-policy: gafaspot-user-ldap`  
ldap-group-policy is the name of a Vault policy attached to tokens created with LDAP Auth Method. When Gafaspot uses the LDAP Auth Method to verify its users, Vault requests over LDAP
* if the user credentials are valid at all and
* in case they are, to which groups the user belongs to.
Depending on the group, Vault associates preconfigured policies to the user and returns the policy names to Gafaspot. Based on this policy name, gafaspot decides, whether the user is allowed to use gafaspot or not.  
For more information about how to configure the LDAP Auth Method correctly, see the instructions about [LDAP Auth Method](doc/auth_ldap.md)
___
`environments:`  
The end of Gafaspot config is for describing the composition of the different environments, which you want to manage with Gafaspot. Therefore, give a list of all environments at the first level like this:

    environments:

        demo0:
            ...

        demo1:
            ...

        demo2:
            ...

        ...
The environment's names are only allowed to contain **lowercase** ASCII letters, numbers and underscores. Don't use uppercase letters and blanks!  
Each environment has further attributes: 

        demo0:
            show-name: DEMO 0
            description: "Some description for DEMO 0;
                          can use multiple lines and
                          HTML tags <br> for formatting."
            secrets-engines:
                ...
As you can see, you are able to provide an attribute `show-name` which is allowed to contain every character. This name will be displayed in web interface. Additionally, the web interface shows every instruction you write into `description`. Use HTML syntax for formatting. For example, you can include hyperlinks. Consider to explain in detail, which components are inside the environment, which credentials are to expect from the Secret Engines, and how the credentials map to the environments. `show-name` and `description` are optional.

Finally, you need to list all Secrets Engines at the third level. Therefore, enable as many Secrets Engines in Vault as you need to perform credential changing for all devices in your environment. In parallel, enable one KV Secrets Engine for each credential-changing secrets engine. You have to enable the Secrets Engines at following paths:

    operate/<environment_name>/<secrets_engine_name>    => Some Secrets Engine offering new credentials
    store/<environment_name>/<secrets_engine_name>      => KV Secrets Engine which stores the credentials for the other Secrets Engine
In this example, environment_names would be `demo0`, `demo1` and so on. For enabling Secrets Engines at the right paths, read the [General Instructions about Secrets Engines](secengs_general.md).

Defining the Secrets Engines looks like this:

            ...
            secrets-engines:
                - name: NetApp
                  type: ontap
                  role: gafaspot
                
                - name: ActiveDirectory
                  type: ad
                  role: gafaspot

                - ...
The Secrets Engine's name may only contain ASCII letters, numbers and underscores. Anyway, try to choose a descriptive name, as this name will be displayed in web interface when user request credentials. As described in [Secrets Engines General](secengs_general.md), the name is the last part of the path, under which you enable the Secrets Engine in Gafaspot.  
`type` is one of:
* ad
* ssh
* database
* ontap

You don't need to explicitly mention KV Secrets Engines in the config file, as they are always related to another Secrets Engine.

`role` is the name of the role you configure with the respective Secrets Engine. How you create the role is described in the respective Instructions about the Secrets Engine type.