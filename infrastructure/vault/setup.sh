#!/bin/bash

# 1. Вмикаємо аутентифікацію через Kubernetes
# Це дозволяє Vault перевіряти JWT токени подів.
echo "🔓 Enabling Kubernetes Auth..."
kubectl exec vault-0 -- vault auth enable kubernetes

# 2. Налаштовуємо Vault, щоб він знав, де знаходиться API кластера
# Vault сидить всередині кластера, тому він звертається до нього через внутрішні змінні.
echo "⚙️ Configuring Kubernetes Auth..."
kubectl exec vault-0 -- sh -c '
  vault write auth/kubernetes/config \
    kubernetes_host="https://$KUBERNETES_PORT_443_TCP_ADDR:443"
'

# 3. Записуємо секрет (пароль до бази даних)
# secret/todo-app/config - це "шухлядка"
# username=postgres, password=megasecret - це вміст
echo "🤫 Writing secrets..."
kubectl exec vault-0 -- vault kv put secret/todo-app/config username="postgres" password="megasecret"

# 4. Завантажуємо файл політики (який ми створили в Кроці 2) всередину пода Vault
echo "📜 Uploading policy..."
kubectl cp infrastructure/vault/policies/backend-policy.hcl vault-0:/tmp/backend-policy.hcl
# Реєструємо цю політику під назвою 'backend-read'
kubectl exec vault-0 -- vault policy write backend-read /tmp/backend-policy.hcl

# 5. Створюємо РОЛЬ (Role) - це найголовніше!
# Тут ми кажемо: "Якщо приходить ServiceAccount з іменем 'backend-sa',
# дай йому права політики 'backend-read'".
echo "🔗 Binding Role..."
kubectl exec vault-0 -- vault write auth/kubernetes/role/backend-role \
    bound_service_account_names=backend-sa \
    bound_service_account_namespaces=default \
    policies=backend-read \
    ttl=24h

echo "✅ Vault successfully configured!"