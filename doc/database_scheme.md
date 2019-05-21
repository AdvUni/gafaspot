# SQLite Database Scheme

Gafaspot uses a simple SQLite Database to store some information persistently. The database location is determined in the [config file](./config_explanation.md). If a database does not exist yet at the given location, Gafaspot will create one at startup. All necessary tables will be created automatically, so you will not have to do any database configuration at all.

Anyway, it might be good to have an overview over the database's contents, so the following graphic shows the database scheme used by Gafaspot:

![database scheme](img/db_scheme.svg)

## Tables
The table `reservations` stores all information about reservations created by users through the web interface. Each reservation has a `status` which is one of the following:
* `upcoming`
* `active`
* `expired`
* `error`

Gafaspot scans the reservations regularly, compares their `start`, `end` and `delete_on` columns with the current point in time, decides whether any actions are necessary, and eventually changes their status accordingly.

The table `environments` gets recreated each time Gafaspot starts to apply possible changes made in the config file. `env_plain_name` and `env_nice_name` correspond to the different identifiers for environments given in the configuration.

The table `users` is for storing public SSH keys which are uploaded by users through the web interface. SSH keys are needed to perform reservations for environments with the SSH Secrets Engine. Entries in table `users` will not be created unless a user uploads a key. Users without a key can still create reservations for environments which do not use the SSH Secrets Engine.

## Relations
The column names *`username`* and *`env_plain_name`* in the table `reservations` are italic in the database scheme and therefore marked as foreign keys of the other tables. Therefore, there are `1:n` relations between these tables. However, those are not real database relations. There are legitimate reasons why the corresponding user or environment entry for a reservation may not exists within the database. This is, for example, the case if a user has not uploaded an SSH key yet. Furthermore, it can happen that after a restart of Gafaspot some environments disappear from database because the configuration has changed. This should have no effect on expired reservations. To make such cases possible, there are no dependencies manifested in the database. Instead, keeping the tables consistent is the job of Gafaspot itself.

## Database manipulations
There are a few direct database manipulations you might want to perform as administrator of gafaspot to control the flow of reservations:
* You can always **delete upcoming reservations** from the database. This will cancel the reservation without causing further trouble.
* You can **change an active reservation's end time** if you want to shorten or extend a reservation which is already active. If the environment concerned by this reservation contains an SSH Secrets Engine, Gafaspot will not be able to adopt these changes to the created SSH certificates. So keep in mind, that the validity period of SSH credentials will not comply with the reservation period anymore if you perform such an operation.
* You **must not delete active reservations** since Gafaspot will not be able to end them properly anymore.
* Reservations with status `expired` or `error` may be deleted any time.


---
*Go back to [table of contents](README.md)...*
