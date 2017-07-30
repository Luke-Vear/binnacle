#!/bin/bash

set -eu

VERSION=${1}

APP_NAME="binnacle"
SRC_PATH="github.com/Luke-Vear/${APP_NAME}"
PORT=8080

echo "Building app using temp container"
docker run --rm -v "${PWD}":/usr/local/go/src/${SRC_PATH} -w /usr/local/go/src/${SRC_PATH} golang:latest \
  /bin/bash -c 'update-ca-certificates
    cp /etc/ssl/certs/ca-certificates.crt .
    go get
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main'

echo "Building app container"
docker build -t ${APP_NAME}:${VERSION} -f Dockerfile .

#echo "Running container on port ${PORT}"
#docker run -d -p ${PORT}:${PORT} ${APP_NAME}:latest /main
#firefox localhost:${PORT}

RM_LIST=("ca-certificates.crt" "${APP_NAME}" "main")

for del in ${RM_LIST[@]}; do

  if [[ -f ${del} ]]; then
    rm -vf ${del}
  fi

done

./publish.sh ${VERSION}
