kind create cluster --config kind-config.yml

helm install rabbitmq oci://registry-1.docker.io/bitnamicharts/rabbitmq

helm install rest-api .
