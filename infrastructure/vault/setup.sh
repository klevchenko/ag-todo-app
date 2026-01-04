#!/bin/bash

# Оновлений скрипт для namespace: todo-app

echo "⏳ Waiting for Vault to be ready..."
# Чекаємо поки под запуститься, щоб скрипт не впав
kubectl wait --for=condition=Ready pod/vault-0 -n todo-app --timeout=60s

echo "🔓 Enabling Kubernetes Auth..."
# Звертаємось до vault-0 саме в todo-app
kubectl exec vault-0 -n todo-app -- vault auth enable kubernetes

echo "⚙️ Configuring Kubernetes Auth..."
kubectl exec vault-0 -n todo-app -- sh -c '
  vault write auth/kubernetes/config \
    kubernetes_host="https://$KUBERNETES_PORT_443_TCP_ADDR:443"
'

echo "🤫 Writing secrets..."
# Записуємо пароль до бази
kubectl exec vault-0 -n todo-app -- vault kv put secret/todo-app/config username="postgres" password="megasecret"

echo "📜 Uploading policy..."
kubectl cp infrastructure/vault/policies/backend-policy.hcl todo-app/vault-0:/tmp/backend-policy.hcl
kubectl exec vault-0 -n todo-app -- vault policy write backend-read /tmp/backend-policy.hcl

echo "🔗 Binding Role..."
# УВАГА: Тут дві зміни!
# 1. bound_service_account_namespaces=todo-app (хто приходить)
# 2. Ми виконуємо це всередині Vault, який теж живе в todo-app
kubectl exec vault-0 -n todo-app -- vault write auth/kubernetes/role/backend-role \
    bound_service_account_names=backend-sa \
    bound_service_account_namespaces=todo-app \
    policies=backend-read \
    ttl=24h

echo "✅ Vault fully configured in 'todo-app' namespace!"