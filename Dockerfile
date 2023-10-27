FROM golang:alpine as build
RUN apk add --no-cache --update git
ADD . /go/src/app
WORKDIR /go/src/app
RUN go get ./...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -tags timetzdata -o yet main.go

# starting with deno runtime
FROM denoland/deno:alpine
# adding ffmpeg
RUN apk update && apk add && apk add ffmpeg
# adding yet
COPY --from=build /go/src/app/yet /usr/bin/yet

#temporary directory
VOLUME /usr/share/yet
