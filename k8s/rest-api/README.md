helm upgrade rest-api --set image.tag=$(git rev-parse --short HEAD) .
