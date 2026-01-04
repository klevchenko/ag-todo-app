# infrastructure/vault/policies/backend-policy.hcl

# Ми дозволяємо action "read" для шляху secret/data/...
# data - це обов'язковий префікс для KV-v2 сховища
path "secret/data/todo-app/config" {
  capabilities = ["read"]
}