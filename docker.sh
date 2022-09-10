#! /usr/bin/env sh

docker build -t eyedeekay/soap .
docker rm -f soap
docker run -itd --net=host --name soap --restart=always eyedeekay/soap
