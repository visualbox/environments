# Build image
docker build -t pwntus/visualbox-node:11 -f ./Dockerfile ..

# Push image to repository
docker push pwntus/visualbox-node:11
