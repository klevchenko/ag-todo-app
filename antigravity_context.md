# 🌌 Project Context: Cloud-Native K8s Todo App

## 1. Project Overview
This is a Fullstack Monorepo application deployed on **Kubernetes (Minikube)**.
It implements a **Microservices Architecture** with a focus on **Zero Trust Security** using HashiCorp Vault.

**Goal:** A "To-Do" list application where Frontend talks to Backend via Reverse Proxy, and Backend retrieves secrets dynamically from Vault.

## 2. Tech Stack
* **Orchestration:** Kubernetes (Minikube).
* **Frontend:** Vue.js 3 (Composition API) + Vite.
* **Frontend Serving:** Nginx (runs in the same pod/container as Vue artifact, acts as Reverse Proxy).
* **Backend:** Go (Golang) 1.24 + Gin Framework + GORM.
* **Database:** PostgreSQL (Bitnami Helm Chart).
* **Secret Management:** HashiCorp Vault (Sidecar Injection Pattern).
* **Ingress:** Nginx Ingress Controller (`todo.local`).
* **Local Env:** Docker Desktop (Mac) + Minikube Tunnel.

## 3. Architecture & Data Flow
1.  **User** accesses `http://todo.local`.
2.  **Ingress Controller** routes traffic to `frontend-svc`.
3.  **Frontend Pod (Nginx)** serves static files.
    * Requests to `/api/*` are proxied by Nginx to `backend-svc`.
4.  **Backend Pod (Go)** receives the request.
    * **Secrets:** Upon startup, `vault-agent-injector` (Sidecar) authenticates with K8s, retrieves secrets from Vault, and writes them to `/vault/secrets/config` (shared volume).
    * **Go App:** Reads config from that file and connects to DB.
5.  **PostgreSQL** executes the query.

## 4. Directory Structure
```text
.
├── apps
│   ├── backend         # Go Code (main.go, Dockerfile)
│   └── frontend        # Vue Code (src/, nginx.conf, Dockerfile)
├── infrastructure
│   ├── k8s-base        # K8s Manifests (Deployment, Service, Ingress)
│   │   ├── backend/
│   │   ├── frontend/
│   │   └── ingress.yaml
│   └── vault           # Vault Setup Scripts & Policies
└── README.md
5. Key Infrastructure Details
Network
Domain: todo.local (mapped to 127.0.0.1 via /etc/hosts + minikube tunnel).

Backend Port: 8080 (Container), 80 (Service).

Frontend Port: 80 (Container/Nginx), 80 (Service).

Vault Configuration
Auth Method: Kubernetes Auth.

Injection: Done via Annotations in backend/deployment.yaml.

Template: Vault Agent creates a .env formatted file at /vault/secrets/config.

Secrets Path: secret/data/todo-app/config.

Docker & Minikube
We use the Minikube Docker daemon: eval $(minikube docker-env).

Image Pull Policy is set to Never to use local images.

6. Common Operational Commands
Build & Deploy
Bash

# 1. Switch Docker Context
eval $(minikube docker-env)

# 2. Build Images (Force No Cache if code changed)
docker build --no-cache -t todo-backend:vX -f apps/backend/Dockerfile apps/backend/
docker build -t todo-frontend:vX -f apps/frontend/Dockerfile apps/frontend/

# 3. Update K8s (after changing image tag in deployment.yaml)
kubectl apply -f infrastructure/k8s-base/backend/
kubectl apply -f infrastructure/k8s-base/frontend/

# 4. Restart Pods (to pick up new code/secrets)
kubectl delete pod -l app=backend -n todo-app
Access & Debugging
Bash

# Access Vault UI
kubectl port-forward svc/vault 8200:8200 -n todo-app

# Logs Backend
kubectl logs -l app=backend -n todo-app -c backend -f

# Check Ingress IP
kubectl get ingress -n todo-app