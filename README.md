install k3s

helm --namespace metallb-system install --create-namespace metallb metallb/metallb
kubectl apply -f metallb_config.yml
kubectl apply -f metallb_announce.yml

helm install rabbitmq oci://registry-1.docker.io/bitnamicharts/rabbitmq

helm install rest-api .

bazel run //src/rest-api/cmd:push -- --tag "$(git rev-parse --short HEAD)"
helm upgrade rest-api --set image.tag=$(git rev-parse --short HEAD) k8s/rest-api
