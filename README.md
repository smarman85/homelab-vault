# homelab-vault

## Initial set up
```bash
$ docker-compose up -d --build

$ docker exec -it vault bash
vault@3134a23f2cd1:/$ consul keygen
wLi63ZredCtC9oSo1mPp3s7TP+ghf144FhYpGfq2K3Y=   ### don't use this one obviously. 
## add these too:
vim consul/config/consul-config.json
vim vault/config/consul-config.json

$ docker-compose up -d --build

$ docker exec -it consul bash
root@e47c12f5e43d:/# consul acl bootstrap
AccessorID:       XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
SecretID:         XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
Description:      Bootstrap Token (Global Management)
Local:            false
Create Time:      2020-01-10 17:46:38.6665438 +0000 UTC
Policies:
   00000000-0000-0000-0000-000000000001 - global-management

root@e47c12f5e43d:/# export CONSUL_HTTP_TOKEN=< SecretID >
root@e47c12f5e43d:/# consul acl policy create -name vault -rules @/consul/vault-policy.hcl
ID:           f5c286ab-ddd2-ee06-7ef7-9fa60767c92a
Name:         vault
Description:
Datacenters:
Rules:
{
  "key_prefix": {
    "vault/": {
      "policy": "write"
....

root@e47c12f5e43d:/# consul acl token create -description "vault agent token" -policy-name vault
AccessorID:       xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
SecretID:         xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
Description:      vault agent token
Local:            false
Create Time:      2020-01-10 17:49:00.0499179 +0000 UTC
Policies:
   f5c286ab-ddd2-ee06-7ef7-9fa60767c92a - vault

root@e47c12f5e43d:/# consul acl set-agent-token agent < SecretID from above >
ACL token "agent" set successfully
exit

# add the token to the vault config
vim vault/config/vault-config.json

$ docker-compose up -d --build
```

## Setting up vault
```bash
$ docker exec -it vault bash
# for ease we will just use 1 unseal key

!!!!!!!!!!!! Note the unseal/root key can not be recovered after this. Keep track of them
vault@1284d3e54a14:/$ vault operator init -n 1 -t 1
Unseal Key 1: XXXXXXXXXXXXXXXXXXXXXXXXXXXXX

Initial Root Token: XXXXXXXXXXXXXXXXXXXXXXX

Vault initialized with 1 key shares and a key threshold of 1. Please securely
distribute the key shares printed above. When the Vault is re-sealed,
restarted, or stopped, you must supply at least 1 of these keys to unseal it
before it can start servicing requests.

Vault does not store the generated master key. Without at least 1 key to
reconstruct the master key, Vault will remain permanently sealed!

It is possible to generate new unseal keys, provided you have a quorum of
existing unseal keys shares. See "vault operator rekey" for more information.


# unseal vault with unseal key from above
vault@1284d3e54a14:/$ vault operator unseal
Unseal Key (will be hidden):
Key                    Value
---                    -----
Seal Type              shamir
Initialized            true
Sealed                 false
Total Shares           1
Threshold              1
Version                1.3.1
Cluster Name           vault-cluster-3df88405
Cluster ID             dd481f31-119b-bb62-c59a-a732ecda3443
HA Enabled             true
HA Cluster             n/a
HA Mode                standby
Active Node Address    <none>

# login to vault with root
vault@1284d3e54a14:/$ vault login
Token (will be hidden):
Success! You are now authenticated. The token information displayed below
is already stored in the token helper. You do NOT need to run "vault login"
again. Future Vault requests will automatically use this token.

Key                  Value
---                  -----
token                XXXXXXXXXXXXXXXXXXXXXX
token_accessor       XXXXXXXXXXXXXXXXXXXXXX
token_duration       ∞
token_renewable      false
token_policies       ["root"]
identity_policies    []
policies             ["root"]
vault@1284d3e54a14:/$

# enable your first secrets engine
vault@1284d3e54a14:/$ vault secrets enable -version=2 -path=secret/ kv
Success! Enabled the kv secrets engine at: secret/
vault@1284d3e54a14:/$ vault kv put secret/test
Must supply data
vault@1284d3e54a14:/$ vault kv put secret/test message=hi
Key              Value
---              -----
created_time     2020-01-10T18:16:26.0016551Z
deletion_time    n/a
destroyed        false
version          1
vault@1284d3e54a14:/$ vault kv get secret/test message=hi
Too many arguments (expected 1, got 2)
vault@1284d3e54a14:/$ vault kv get secret/test
====== Metadata ======
Key              Value
---              -----
created_time     2020-01-10T18:16:26.0016551Z
deletion_time    n/a
destroyed        false
version          1

===== Data =====
Key        Value
---        -----
message    hi

```

## Further Info/Learning:
[Learn Vault](https://learn.hashicorp.com/vault)
