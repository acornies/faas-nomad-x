#!/bin/bash

echo 'Waiting for vault...'
while true
do
  START=`docker logs dev-vault 2>&1 | grep "post-unseal setup complete"`
  if [ -n "$START" ]; then
    break
  else
    sleep 2
  fi
done

POLICY_NAME=openfaas
TOKEN=vagrant
VAULT_URL=http://127.0.0.1:8200

export VAULT_ADDR=${VAULT_URL}
export VAULT_TOKEN=${TOKEN}

vault secrets disable secret/
vault secrets enable -path=secret -version=1 kv
vault secrets enable -version=2 kv

vault auth enable approle

curl https://nomadproject.io/data/vault/nomad-server-policy.hcl -O -s -L
curl https://nomadproject.io/data/vault/nomad-cluster-role.json -O -s -L
vault policy write nomad-server nomad-server-policy.hcl
vault write /auth/token/roles/nomad-cluster @nomad-cluster-role.json

# openfaas vault policy
vault policy write ${POLICY_NAME} /vagrant/contrib/vault/policy.hcl

# create basic auth secrets
curl -i --header "X-Vault-Token: ${TOKEN}" \
  --request POST \
  --data '{"username":"admin", "password":"vagrant"}' \
  ${VAULT_URL}/v1/secret/openfaas/auth/credentials

# create approle openfaas
curl -i \
  --header "X-Vault-Token: ${TOKEN}" \
  --request POST \
  --data '{"policies": "openfaas", "period": "5m"}' \
  ${VAULT_URL}/v1/auth/approle/role/${POLICY_NAME}