# Path operate/ holds all credential changing secrets engines
path "operate/*" {
  capabilities = ["create", "read", "update", "delete"]
}

# Path store/ holds all KV secrets engines which store credentials from secrets engines at path operate/
path "store/*" {
  capabilities = ["create", "read", "update", "delete"]
}

# Gafaspot uses this path to tune the default and max ttl for leases created by Secrets Engines
path "sys/mounts/operate/*" {
  capabilities = ["update"]
}

# Gafaspot needs orphan tokens with individual life spans for starting
# reservations to ensure that leases are not revoked to early
path "auth/token/create-orphan" {
  capabilities = ["update"]
}

# Gafaspot uses this path to tune the max ttl for orphan tokens
path "sys/mounts/auth/token/tune" {
  capabilities = ["update"]
}
