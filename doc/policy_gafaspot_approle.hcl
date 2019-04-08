# Path operate/ holds all credential changing secrets engines
path "operate/*" {
  capabilities = ["create", "read", "update", "delete"]
}

# Path store/ holds all KV secrets engines which store credentials from secrets engines at path operate/
path "store/*" {
  capabilities = ["create", "read", "update", "delete"]
}

# Gafaspot uses this path to overwrite the default ttl for leases
path "sys/mounts/operate/*" {
  capabilities = ["create", "read", "update", "delete"]
}

# Access to this path is needed to revoke leases created from secrets engines at path operate/
path "sys/leases/revoke-prefix/operate/*" {
  capabilities = ["create", "read", "update", "delete"]
}