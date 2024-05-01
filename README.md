# yet
yet is a minimalist YouTube video, playlist, channel downloader. 

Engineering design constraints lead to a simple application and code base. yet is built on top of `github.com/boggydigital/yt_urls`, similarly compact module, that provides low-level helpers to get and work with YouTube data.

yet can be run as a [CLI tool](#using-yet-as-a-cli-tool) or as a [Docker service](#using-yet-as-a-docker-service).

## Using yet as a CLI tool

To download a video use the following command:

```shell
yet download <video-id> [, <another-video-id>...] | <playlist-id> [, <another-playlist-id>...]
```

yet supports individual video-ids or playlist-ids as args.

Internally yet downloads videos using a list of video-ids, so any playlist-ids are expanded into video-ids for download.

### What is video-id?

Video-id is YouTube's video identifier. You can get it from a video URL : `https://www.youtube.com/watch?v=video-id`.

### What is playlist-id?

Playlist-id is YouTube's videos list identifier. You can get it from a list URL: `https://www.youtube.com/watch?v=video-id&list=playlist-id`.

### Setting up directories

yet requires several directories for operations - see [`directories-example.txt`](https://github.com/boggydigital/yet/blob/main/directories-example.txt) and replace default paths with desired values. Rename and place `directories.txt` in the same folder where `yet` binary is located.

## Using yet as a Docker service

The recommended way to install yet is with docker-compose:

- create `compose.yml` file (this minimal example omits common settings like network, restart, etc):

```yaml
version: '3'
services:
  yet:
    container_name: yet
    image: ghcr.io/boggydigital/yet:latest
    volumes:
      # backups
      - /docker/yet/backups:/usr/share/yet/backups
      # input
      - /docker/yet/input:/usr/share/yet/input
      # metadata
      - /docker/yet/metadata:/usr/share/yet/metadata
      # videos
      - /docker/yet/videos:/usr/share/yet/videos
      # posters
      - /docker/yet/posters:/usr/share/yet/posters
      # captions
      - /docker/yet/captions:/usr/share/yet/captions
      # players
      - /docker/yet/players:/usr/share/yet/players
      # sharing timezone from the host
      - /etc/localtime:/etc/localtime:ro 
```
- (move it to location of your choice, e.g. /docker/yet or remote server or anywhere else)
- while in the directory with that config - pull the image with `docker compose pull`
- run yet with `docker compose run yet yet download <video-id>` to download a video to `videos` directory

yet Docker service is a work in progress and will be improved in the future to add more functionality. 

## Advanced scenarios

Despite supporting only video-ids and playlist-ids, yet also (implicitly) supports channel and user videos, when they can be expressed as playlists. 

All advanced scenarios are detailed below. 

### Downloading all channel and user videos

Most channels and user pages contain playlist links to all uploads. For example: on a channel page look for "Uploads" section and a "PLAY ALL" link - this link can be used as a playlist-id.

### Downloading videos that require authentication

Some YouTube videos require users to sign in - e.g. paid-membership only videos.

yet uses `coost` to persist session cookies and it's possible to reuse existing browser YouTube cookies to get access to videos that require authentication. 

Please refer to [the coost README](https://github.com/boggydigital/coost#copying-session-cookies-from-an-existing-browser-session) for the step by step guide on copying YouTube session cookies.

For yet you need to create or edit `cookies.txt` file in the yet input directory and add `youtube.com` host sessions cookies. It should look like this:

```text
youtube.com
  cookie-header=<paste-youtube-session-cookie-header-from-your-browser-here>
```

### Testing yet dependencies

yet can be enhanced with several improvements when used in an environment with `ffmpeg`, `node` or `deno` - better video and audio quality, faster downloads, etc.

yet Docker image contains `ffmpeg` and `deno` and provides the best quality and fastest downloads out of the box - no other action needed.

yet CLI usage needs external binaries available and it's recommended to test the dependencies before attempting to setup them. yet provides a command for that:

```shell
yet test-dependencies
```
This will output something like this:
```text
testing dependencies... 
 ffmpeg /usr/local/bin/ffmpeg 
 node /usr/local/bin/node 
 deno not found 
```

### Specifying ffmpeg binary path to get the best quality video/audio

yet Docker image should include `ffmpeg`, so this section will help you setting up (or not) `ffmpeg` for CLI tool usage. 

Please note that `ffmpeg` is NOT required for yet to function - yet was designed to function without any external dependencies out of the box. When yet cannot locate a working `ffmpeg` binary, it'll download a mobile version of the video that'll contain video and audio in one file. Typically that means 720p videos / medium quality sound and is a YouTube limitation, not yet.

However, if an external dependency is not a problem for your use-case - you can progressively enhance yet with `ffmpeg`. By default, yet will attempt to locate `ffmpeg` binary on the system. In most cases that's sufficient and assuming you have `ffmpeg` installed - you don't need to do anything special to get better quality video/audio.

If you'd prefer to specify `ffmpeg` binary location manually, set `YET_FFMPEG_CMD` environment variable to the full path of `ffmpeg` binary (e.g. `/opt/homebrew/bin/ffmpeg` for Homebrew installation on macOS).

### Enabling faster downloads with a JavaScript engine

YouTube implements measures to restrict download speed, unless download client passes a challenge. [This issue](https://github.com/ytdl-org/youtube-dl/issues/29326#issuecomment-894619419) goes into more details - check it out if you want to know more about that restriction.

In order to unlock faster downloads yet can extract decoding code from YouTube video page and run it for you. This but requires a JavaScript engine (Node.js or Deno) that would run that decoding code.

Below you will find details on how to enable each option depending on your needs and available software. There is of course a third option - do nothing and live with slower yet download speeds.

#### Using a JavaScript engine to run decoding code automatically

If an external dependency is not a problem for your use-case or you already have Node.js installed - you can progressively enhance yet with `node` or `deno`. By default, yet will attempt to locate `node` and `deno` binaries on the system. In most cases that's sufficient and assuming you have Node.js or Deno installed - you don't need to do anything special to unlock faster download speed. 

If you'd prefer to specify `node` or `deno` binary location manually, set `YET_NODE_CMD` environment variable for `node` or `YET_DENO_CMD` for `deno` to the full path of the binary (e.g. `/usr/local/bin/node` for default Node.js installation on macOS).

However you specify location of the `node` or `deno` binary - upon encountering an encoded parameter, yet will download decoding code and create `decoder.js` file in the working directory and use an available JavaScript runtime binary to run it and get the decoded value - this is completely automatic and doesn't require user input.

## Building for another OS

Go allows you to build binary for another OS like this (using Linux and AMD64 as an example):

```shell
env GOOS=linux GOARCH=amd64 go build -o yet
```

## Privacy

yet doesn't collect any data whatsoever. Whatever you do with yet stays on your machine. 

Your Internet connection is only used to download YouTube metadata and videos and nothing is ever uploaded anywhere. 

If you've provided any `youtube.com` cookies - they're transmitted as part of requests to get YouTube data, exactly the same way your browser would send them. You can delete `cookies.txt` that you've created at any point with no impact to ability to download publicly available videos (you won't be able to download any YouTube videos that require authorization). 
