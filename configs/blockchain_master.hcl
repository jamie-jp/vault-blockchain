path "blockchain/keys/create" {
    capabilities = [ "create" ]
}

path "blockchain/keys/+/sign" {
    capabilities = [ "read" ]
}
