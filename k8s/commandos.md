# Criar cluester kubernete
```bash
k3d cluster create nlw-unite --servers 2
```
# Comandos kubectl
## Criar namespace
```bash
kubectl create ns nlw
```
## Aplicar as configurações
```bash
kubectl apply -f /path/to/dir -n nlw
```
## Pegar todos os deployments
```bash
kubectl get deployments -n nlw
```
## Pegar todos os pods
```bash
kubectl get pods -n nlw
```
