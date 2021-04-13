path "blockchain/create" {
    capabilities = [ "create" ]
}

path "blockchain/accounts/+/sign" {
    capabilities = [ "read" ]
}
