# 🚀 Cloud-Native Todo App (Monorepo)

A full-stack, microservices-based Todo application deployed on Kubernetes.
This project demonstrates a production-ready architecture using **Go**, **Vue.js**, **PostgreSQL**, **HashiCorp Vault** for secret management, and **Nginx Ingress** for routing.

![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Vue.js](https://img.shields.io/badge/vuejs-%2335495e.svg?style=for-the-badge&logo=vuedotjs&logoColor=%234FC08D)
![Vault](https://img.shields.io/badge/vault-%236E7681.svg?style=for-the-badge&logo=vault&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)

## 🏗 Architecture

The application follows a **Zero Trust** security model where the backend retrieves database credentials dynamically from HashiCorp Vault via a Sidecar Agent injection.

* **Frontend:** Vue.js 3 (Vite) + Nginx (Reverse Proxy).
* **Backend:** Go (Gin Framework) + GORM.
* **Database:** PostgreSQL (Running in K8s).
* **Security:** HashiCorp Vault (Secrets Management & Agent Injection).
* **Routing:** Nginx Ingress Controller.

---

## 📂 Project Structure

```text
.
├── apps
│   ├── backend         # Go API (Gin + GORM)
│   └── frontend        # Vue.js App + Nginx Config
├── infrastructure
│   ├── k8s-base        # Kubernetes Manifests (Deployment, SVC, Ingress)
│   └── vault           # Setup scripts & Policies
└── README.md
🛠 Prerequisites
Ensure you have the following tools installed locally:

Docker Desktop (or Engine)

Minikube (Kubernetes Cluster)

Kubectl (K8s CLI)

Helm (Package Manager)

🚀 Getting Started
1. Start Minikube
Initialize the cluster and enable the Ingress addon.

Bash

minikube start
minikube addons enable ingress
2. Namespace & Context
Create the dedicated namespace and switch context to avoid typing -n todo-app every time.

Bash

kubectl create namespace todo-app
kubectl config set-context --current --namespace=todo-app
3. Install Infrastructure (Vault & Postgres)
We use Helm to install the base infrastructure.

Bash

# Install Vault (Dev mode)
helm install vault hashicorp/vault \
  --set "server.dev.enabled=true" \
  --set "ui.enabled=true" \
  --namespace todo-app

# Install Postgres
helm install postgres bitnami/postgresql \
  --set auth.postgresPassword=megasecret \
  --set primary.persistence.enabled=false \
  --namespace todo-app
4. Configure Vault (Auto-Script)
Run the setup script to enable Kubernetes Auth, create roles, and write initial secrets.

Bash

chmod +x infrastructure/vault/setup.sh
./infrastructure/vault/setup.sh
5. Build Application Images
Build Docker images directly inside Minikube's environment.

Bash

eval $(minikube docker-env)

# Build Backend
docker build -t todo-backend:v5 -f apps/backend/Dockerfile apps/backend/

# Build Frontend
docker build -t todo-frontend:v1 -f apps/frontend/Dockerfile apps/frontend/
6. Deploy Applications
Apply Kubernetes manifests for Backend, Frontend, and Ingress.

Bash

kubectl apply -f infrastructure/k8s-base/backend/
kubectl apply -f infrastructure/k8s-base/frontend/
kubectl apply -f infrastructure/k8s-base/ingress.yaml
🌐 Accessing the App
1. Start Tunnel (Required for macOS users)
Keep this command running in a separate terminal window to expose the LoadBalancer IP.

Bash

minikube tunnel
2. Configure DNS
Add the local domain to your /etc/hosts file.

Bash

# Get the Ingress IP (usually 127.0.0.1 with tunnel)
kubectl get ingress

# Edit hosts file
sudo nano /etc/hosts
Add this line:

Plaintext

127.0.0.1 todo.local
3. Open in Browser
Visit http://todo.local

🔐 Managing Secrets (Vault)
To add new secrets dynamically without rebuilding the image:

Access Vault UI:

Bash

kubectl port-forward svc/vault 8200:8200
URL: http://localhost:8200

Token: root

Add Secret: Navigate to secret/todo-app/config and add Key/Value pairs.

Update Deployment: Ensure the key is added to the agent-inject-template in backend/deployment.yaml.

Restart Pod:

Bash

kubectl delete pod -l app=backend
💻 Development Workflow
If you change code (e.g., in apps/backend/main.go):

Rebuild Image:

Bash

eval $(minikube docker-env)
docker build --no-cache -t todo-backend:v6 -f apps/backend/Dockerfile apps/backend/
Update Deployment: Change image version to v6 in deployment.yaml.

Apply & Restart:

Bash

kubectl apply -f infrastructure/k8s-base/backend/deployment.yaml
🛑 Cleanup
To stop the cluster and free up resources:

Bash

# Stop Minikube (saves state)
minikube stop

# Delete Cluster (fresh start next time)
minikube delete

---
