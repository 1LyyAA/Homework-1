
# minikube start

# Install Istio
istioctl install --set profile=demo --set "values.global.proxy.resources.requests.cpu=10m" --set "values.global.proxy.resources.requests.memory=100Mi" -y

# Enable Istio injection in default namespace
kubectl label namespace default istio-injection=enabled

docker build -t example-app:latest .
minikube image load example-app:latest

# kubectl apply -f pod.yaml

kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/daemonset.yaml
kubectl apply -f k8s/cronjob.yaml


kubectl apply -f k8s/gateway.yaml
kubectl apply -f k8s/virtualservice.yaml
kubectl apply -f k8s/destinationrule.yaml

echo "Waiting for Deployment to be ready..."
kubectl rollout status deployment/app-deployment

echo "Waiting for DaemonSet to be ready..."
kubectl rollout status daemonset/log-agent

echo "=== Deployment finished successfully ==="

# minikube stop
