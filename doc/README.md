# Documentation - Table of contents

The `doc/` folder contains following instructions to help you running Gafaspot:

1. [Vault Setup](vault_setup.md)
    1. [LDAP Auth Method](auth_ldap.md)
    1. [AppRole Auth Method](auth_approle.md)
    1. [Secrets Engines](secengs_general.md)
        * [Active Directory Secrets Engine](secengs_ad.md)
        * [SSH Secrets Engine (Signed Certificates)](secengs_ssh.md)
        * [Database Secrets Engine](secengs_database.md)
        * [Ontap Secrets Engine](secengs_ontap.md)
1. [Explanations for Gafaspot Configuration](config_explanation.md)
1. [SQLite Database Scheme](database_scheme.md)

In addition, this folder contains a subdirectory [`json_payload/`](json_payload) containing several JSON snippets which are used in the Vault setup guide. So, you can navigate into this directory and run the proposed commands and the necessary payload will already be in place.

Finally, the `doc/` folder includes a file called [`policy_gafaspot_approle.hcl`](policy_gafaspot_approle.hcl), which is not JSON payload, but also used in the Vault setup guide.