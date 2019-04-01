# Secrets Engines General

Vault's [Secrets Engines](https://www.vaultproject.io/docs/secrets/) may work very differently. Gafaspot is able two deal with following Secrets Engines:
* Secrets Engines which are able to change credentials:
    * **Active Directory** Secrets Engine
    * **Database** Secrets Engine
    * **SSH** Secrets Engine (in mode Signed Certificates)
    * **Ontap** Secrets Engine (not an official Secrets Engine from Hashicorp)
* Secrets Engine to store credentials retrieved from other Secrets Engines: **KV** Secrets Engine (Version 1)

Many of the credential-changing Secrets Engines are able to serve exactly one device -- one Active Directory Secrets Engine is able to change passwords for accounts at one Active Directory. So, you'll need to enable one Active Directory Secrets Engine per Active Directory Instance in your environments. But the SSH Secrets Engine for example is not bound to any device; its signed certificates are valid for as many ssh machines as you would like to. As it is the point of Gafaspot to grant access for different environments independently, you probably would like to enable one (if not none) SSH Secrets Engine per environment. What this wants to say is, that how many Secrets Engines are enabled for the devices in one environment, depends on the Secrets Engines you need. What is always equal is, that each Secrets Engine is only responsible for one environment.

To structure all the Secrets Engines and to relate them to specific environments, the circumstance is used, that Vault allows to enable Secrets Engines under different paths. So, 