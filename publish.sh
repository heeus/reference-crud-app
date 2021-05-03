#!/bin/bash

publish_testservice() {
    echo "Building test service..."
    binTestService="testservice-linux-amd64-${version}"
    binLTestService="testservice-linux-amd64-latest"
    env GOOS=linux GOARCH=amd64 go build -o ./bin/${binTestService}
    if [ -f ./bin/${binTestService} ]; then
        echo "Uploading test service..."
        env AWS_ACCESS_KEY_ID=${keyId} AWS_SECRET_ACCESS_KEY=${secretKey} aws s3 cp ./bin/${binTestService} s3://${bucket}/testservice/${binTestService}
        env AWS_ACCESS_KEY_ID=${keyId} AWS_SECRET_ACCESS_KEY=${secretKey} aws s3 cp s3://${bucket}/testservice/${binTestService} s3://${bucket}/testservice/${binLTestService}
        rm ./bin/${binTestService}
    else        
        echo "Build failed"
    fi        
}

version="$(date '+%Y%m%d%H%M')"

keyId="${HEEUS_RELEASER_AWS_KEY_ID}"
secretKey="${HEEUS_RELEASER_AWS_SECRET_KEY}"
bucket="${HEEUS_S3_BUCKET}"

if [[ -z "${keyId}" || -z "${secretKey}" || -z "${bucket}" ]]; then
    echo "Failed: HEEUS_RELEASER_AWS_KEY_ID, HEEUS_RELEASER_AWS_SECRET_KEY and HEEUS_S3_BUCKET environment variables must be set"
    exit 1
fi

echo "Building version: ${version}..."

mkdir -p ./bin
publish_testservice
echo "done"
