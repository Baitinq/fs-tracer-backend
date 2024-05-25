install k3s
install helmsman

COMMIT_SHA=$(git rev-parse --short HEAD) helmsman --apply -f k8s/helmsman.yml

to deploy, execute the ./deploy.sh script (and have the correct env variables)
