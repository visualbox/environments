# Build image
docker build -t pwntus/visualbox-node-dev:11 -f ./Dockerfile ..

# Push image to repository
docker push pwntus/visualbox-node-dev:11
