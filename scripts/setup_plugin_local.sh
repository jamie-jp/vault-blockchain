#!/bin/bash

function install_plugin {
  make build-local
  export SHASUM256=$(shasum -a 256 "./plugins/vault-blockchain" | cut -d' ' -f1)
  vault write sys/plugins/catalog/blockchain-plugin \
        sha_256="${SHASUM256}" \
        command="vault-blockchain --tls-skip-verify=true"
  vault secrets enable -path=blockchain -description="Wallet" -plugin-name=blockchain-plugin plugin
}

function create_policy {
  vault policy write blockchain_master ./configs/blockchain_master.hcl
}

vault login $ROOT_TOKEN
install_plugin
create_policy

echo 'Root Token: '$ROOT_TOKEN

