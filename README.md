install k3s

helm --namespace metallb-system install --create-namespace metallb metallb/metallb
kubectl apply -f metallb_config.yml
kubectl apply -f metallb_announce.yml

helm install rabbitmq oci://registry-1.docker.io/bitnamicharts/rabbitmq

bazel run //src/rest-api/cmd:push -- --tag "rest-api-$(git rev-parse --short HEAD)"
helm upgrade rest-api --set image.tag="rest-api-$(git rev-parse --short HEAD)" k8s/rest-api

bazel run //src/payload-processor/cmd:push -- --tag "payload-processor-$(git rev-parse --short HEAD)"
helm upgrade payload-processor --set image.tag="payload-processor-$(git rev-parse --short HEAD)" k8s/payload-processor
