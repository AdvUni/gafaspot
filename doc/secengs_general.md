# Secrets Engines General
To enable Gafaspot to change credentials for external devices, Vault's Secrets Engines are used. You define environments through enabling and configuring appropriate Secrets Engines for all devices in your environments.

## Supported Secrets Engines
Vault's [Secrets Engines](https://www.vaultproject.io/docs/secrets/) may work very differently. Gafaspot is able to deal with following Secrets Engines:
* Secrets Engines which are able to change credentials:
    * **Active Directory** Secrets Engine
    * **Database** Secrets Engine
    * **SSH** Secrets Engine (in mode Signed Certificates)
    * **Ontap** Secrets Engine (not an official Secrets Engine from HashiCorp)
* Secrets Engine to store credentials retrieved from other Secrets Engines: **KV** Secrets Engine (Version 1)

## How many Instances to Enable?
Many of the credential-changing Secrets Engines are able to serve exactly one device - one Active Directory Secrets Engine is able to change passwords for accounts at one Active Directory. So, you'll need to enable one Active Directory Secrets Engine per Active Directory Instance in your environments. But the SSH Secrets Engine for example is not bound to any device; its signed certificates are valid for as many ssh machines as you would like to. As it is the point of Gafaspot to grant access for different environments independently, you probably would like to enable one (if not none) SSH Secrets Engine per environment. What this wants to say is, that how many Secrets Engines are enabled for the devices in one environment, depends on the Secrets Engines you need. Always equal is, that each Secrets Engine is only responsible for one environment.

## Secrets Engines Paths
[As you learn in the Getting Started Guide](https://learn.hashicorp.com/vault/getting-started/secrets-engines#enable-a-secrets-engine), each Secrets Engine is enabled under a unique path, similar to a file path. This path system is used to structure all the Secrets Engines and to relate them to specific environments. 

### Paths for Credential-changing Secrets Engines
So, a credential-changing Secrets Engine for Gafaspot has to be enabled at following path:

    operate/<environment_name>/<secrets_engine_name>

For the variables environment_name and secrets_engine_name following conventions must be met:
* environment_name is only allowed to contain **lowercase** ASCII letters, numbers and underscores
* secrets_engines_name is allowed to contain (lowercase and uppercase) ASCII letters, numbers and underscores
* the names environment_name and secrets_engine_name are the same you enter in the gafaspot_config.yaml configuration file (in this documentation, more information about the [config file](config_explanations) will follow)
* obviously, each Secrets Engine path must be unique, so, for one environment no secrets_engine_name may appear twice
* try to give a descriptive name for secrets_engine_name, as it will be directly shown in Gafaspot web interface. Further explanations about which secrets_engine_name means what, can be given with an environment description inside gafaspot_config.yaml. For environment_name, the config file allows you to give an extra name which can contain any kinds of characters to be shown in the web interface

The constant prefix `operate` is to indicate that the Secrets Engines perform any kind of operation, which is capable to change access data for some kind of account. It is used, because there is another kind of Secrets Engine used with Gafaspot: The KV Secrets Engine. 

### Paths for KV Secrets Engines
The [KV (Key-Value) Secrets Engine](https://www.vaultproject.io/docs/secrets/kv/kv-v1.html) is needed, because the other secrets engine are not generally able to remember credentials after creation. So, Gafaspot stores them all into KV Secrets Engines to access them later. For doing this in a consistent way, you have to enable **one KV Secrets Engine per credential-changing Secrets Engine** in Vault. Therefore, you use the same path as for the other Secrets Engine, but replace 'operate' with 'store'.

    operate/<environment_name>/<secrets_engine_name>    => Some Secrets Engine offering new credentials
    store/<environment_name>/<secrets_engine_name>      => KV Secrets Engine which stores the credentials for the other Secrets Engine

## Example
So, a fictive Vault setup may have a Secrets Engines path structure like this:

    Path                              Type
    ----                              ----
    operate/demo0/ActiveDirectory/    ad
    operate/demo0/MySQL/              database
    operate/demo0/SSH/                ssh
    operate/demo1/ActiveDirectory/    ad
    operate/demo1/NetApp/             ontap

    store/demo0/ActiveDirectory/      kv
    store/demo0/MySQL/                kv
    store/demo0/SSH/                  kv
    store/demo1/ActiveDirectory/      kv
    store/demo1/NetApp/               kv

The respective config file for Gafaspot gafaspot_config.yaml, which must follow the same structure, would look like this:

    [...]

    environments:

        demo0:
            show-name: DEMO 0
            description: "this is demo environment 0."
            secret-engines:
              - name: ActiveDirectory
                type: ad
                role: gafaspot

              - name: MySQL
                type: database
                role: gafaspot

              - name: SSH
                type: ssh
                role: gafaspot

        demo1:
            show-name: DEMO 1
            description: "this is demo environment 1."
            secret-engines:
              - name: ActiveDirectory
                type: ad
                role: gafaspot

              - name: NetApp
                type: ontap
                role: gafaspot

"role" is an attribute which you have to configure for each credential-changing secrets engine. More about this and other configuration at the respective pages:

* [Active Directory Secrets Engine](secengs_ad.md)
* [SSH Secrets Engine (Signed Certificates)](secengs_ssh.md)
* [Database Secrets Engine](secengs_database.md)
* [Ontap Secrets Engine](secengs_ontap.md)


---
*Go to [next page](secengs_ad.md)...*  
*Go to [table of contents](README.md)...*