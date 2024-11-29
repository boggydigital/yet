FROM golang:alpine as build
RUN apk add --no-cache --update git
ADD . /go/src/app
WORKDIR /go/src/app
RUN go get ./...
RUN go build \
    -a -tags timetzdata \
    -o yet \
    -ldflags="-s -w -X 'github.com/boggydigital/yet/cli.GitTag=`git describe --tags --abbrev=0`'" \
    main.go

# starting with Node.js runtime
FROM node:alpine
# adding ffmpeg
# RUN apk update && apk add && apk add ffmpeg
# adding yet
COPY --from=build /go/src/app/yet /usr/bin/yet

EXPOSE 2005

# backups
VOLUME /usr/share/yet/backups
# input
VOLUME /usr/share/yet/input
# metadata
VOLUME /usr/share/yet/metadata
# videos
VOLUME /usr/share/yet/videos
# players
VOLUME /usr/share/yet/players

ENTRYPOINT ["/usr/bin/yet"]
CMD ["serve","-port", "2005", "-stderr"]

