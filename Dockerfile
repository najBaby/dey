FROM golang:1.16.4-buster
ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
ENV APP_USER app
ENV APP_HOME /go/src/dey
EXPOSE 9000
CMD ["go", "run"]
docker tag local-image:tagname new-repo:tagname
docker push new-repo:tagname
docker swarm join --token $TOKEN $HOST:$PORT