FROM golang:alpine as build
RUN apk add --no-cache --update git
ADD . /go/src/app
WORKDIR /go/src/app
RUN go get ./...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -tags timetzdata -o yet main.go

FROM denoland/deno:alpine
RUN apt-get -y update && apt-get -y upgrade && apt-get install -y --no-install-recommends ffmpeg
COPY --from=build /go/src/app/yet /usr/bin/yet

#temporary directory
VOLUME /var/tmp