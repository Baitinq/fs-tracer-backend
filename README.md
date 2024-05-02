install k3s

helm --namespace metallb-system install --create-namespace metallb metallb/metallb
kubectl apply -f metallb_config.yml
kubectl apply -f metallb_announce.yml

helm install kafka oci://registry-1.docker.io/bitnamicharts/kafka --set controller.replicaCount=1,controller.livenessProbe.initialDelaySeconds=120

to deploy, execute the ./deploy.sh script
