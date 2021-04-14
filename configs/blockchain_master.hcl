path "blockchain/create" {
    capabilities = [ "create" ]
}

path "blockchain/keys/+/sign" {
    capabilities = [ "read" ]
}
