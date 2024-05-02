install k3s

helm --namespace metallb-system install --create-namespace metallb metallb/metallb
kubectl apply -f metallb_config.yml
kubectl apply -f metallb_announce.yml

helm install rabbitmq oci://registry-1.docker.io/bitnamicharts/rabbitmq

to deploy, execute the ./deploy.sh script
