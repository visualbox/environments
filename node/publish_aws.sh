# Build image
docker build -t visualbox-node .

# STEPS DESCRIBED AT
# https://docs.aws.amazon.com/AmazonECR/latest/userguide/docker-basics.html

# Describe ECR repository
REPOSITORY_URI=$(aws ecr describe-repositories --profile prod --region eu-west-1 | jq -r '.repositories[0].repositoryUri')

# Tag new image with repositoryUri from described repository
docker tag visualbox-node $REPOSITORY_URI

# Generate docker login cmd from AWS CLI and run it
eval $(aws ecr get-login --no-include-email --profile prod --region eu-west-1)

# Push image to repository
docker push $REPOSITORY_URI
