
# minikube start

docker build -t example-app:latest .
minikube image load example-app:latest

# kubectl apply -f pod.yaml

kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/daemonset.yaml
kubectl apply -f k8s/cronjob.yaml



echo "Waiting for Deployment to be ready..."
kubectl rollout status deployment/app-deployment

echo "Waiting for DaemonSet to be ready..."
kubectl rollout status daemonset/log-agent

echo "=== Deployment finished successfully ==="

# minikube stop
