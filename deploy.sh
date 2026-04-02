



minikube start

docker build -t example-app:latest .
minikube image load example-app:latest

kubectl apply -f pod.yaml

kubectl apply -f configmap.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f daemonset.yaml

minikube stop
