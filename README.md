# yet
yet is a minimalist YouTube video and channel downloader. Engineering design constraints lead to a simple application and code base. yet is built on top of `github.com/boggydigital/yt_urls`, similarly opinionated module, that provides low-level helpers to get and work with YouTube data.

## Using yet

```shell
yet <video-id> [, <another-video-id>...] | <playlist-id> [, <another-playlist-id>...]
```

yet supports individual video-ids or playlist-ids as args. Either one can be provided as a full `https://www.youtube.com/...` link.

Internally yet downloads videos using a list of video-ids, so any channel-ids are expanded into video-ids for download.

At the moment, there are no other (explicit) parameters that yet supports. When called without any arguments yet will print help information.

### What is video-id?

Video-id is YouTube's video identifier. You can get it from a video URL : `https://www.youtube.com/watch?v=video-id`. yet can extract video-id from a URL, so you can use either video-id or the full URL `https://www.youtube.com/watch?v=video-id`.

### What is playlist-id?

Playlist-id is YouTube's videos list identifier. You can get it from a list URL: `https://www.youtube.com/watch?v=video-id&list=playlist-id`. Similarly to video-id, yet supports URL containing playlist-id, so you can use a full URL. Please note: URL can contain playlist-id and video-id at the same time and in that case playlist-id will be prioritized over video-id. If that's not desired - make sure to use either URL with just video-id or video-id itself.

## Advanced scenarios

Despite supporting only video-ids and playlist-ids, yet also (implicitly) supports channel and user videos, when they can be expressed as playlists. 

All advanced scenarios are detailed below. 

### Downloading channel and user videos

Most channels and user pages contain playlist links to all uploads. For example: on a channel page look for "Uploads" section and a "PLAY ALL" link - this link can be used as a playlist-id.

### Downloading videos that require authentication

Some YouTube videos require users to sign in - e.g. paid membership only videos.

yet uses coost to persist session cookies and it's possible to reuse existing browser YouTube cookies to get access to videos that require authentication. 

Please refer to [the coost README](https://github.com/boggydigital/coost#copying-session-cookies-from-an-existing-browser-session) for the step by step guide on copying YouTube session cookies.

For yet you need to create or edit `cookies.txt` file in the yet working directory and add `youtube.com` host sessions cookies. It should look like this:

```text
youtube.com
  cookie-header=<paste-youtube-session-cookie-header-from-your-browser-here>
```
### Specifying ffmpeg location to get the best quality video/audio

Please note that `ffmpeg` is NOT required for yet to function - yet was designed to function without any external dependencies out of the box. When yet cannot locate a working `ffmpeg` binary, it'll download a mobile version of the video that'll contain video and audio in one file. Typically that means 720p videos / medium quality sound and is a YouTube limitation, not yet.

However, if an external dependency is not a problem for your use-case - you can progressively enhance yet with `ffmpeg`. By default, yet will attempt to locate `ffmpeg` binary on the system. In most cases that's sufficient and assuming you have `ffmpeg` installed - you don't need to do anything special to get better quality video/audio.

If you'd prefer to specify `ffmpeg` binary location manually, set `YET_FFMPEG_CMD` environment variable to the full path of `ffmpeg` binary (e.g. `/opt/homebrew/bin/ffmpeg` for Homebrew installation on macOS).

### Enabling faster downloads

YouTube implements measures to restrict download speed, unless download client passes a challenge. [This issue](https://github.com/ytdl-org/youtube-dl/issues/29326#issuecomment-894619419) goes into more details - check it out if you want to know more about that restriction.

In order to unlock faster downloads yet can extract decoding code from YouTube video page and run it for you. There are two ways to achieve that. THe first one requires user input, but doesn't require any additional software other than the web browser that has JavaScript engine to run that decoding code (any modern browser would work). The second option doesn't require user input and is completely automatic, but requires Node.js that would run that decoding code.

Below you will find details on how to enable each option depending on your needs and available software. There is of course a third option - do nothing and live with slower yet download speeds.

#### Using web browser to run decoding code

Set `YET_FAST` environment variable to any non-empty value and run yet. For example on macOS you can do the following:

```shell
YET_FAST=1 ./yet <video-id>
```
Upon encountering an encoded parameter, yet will download decoding code and create `decoder.html` file in the working directory and post a `file://` link to it in your Terminal. You would need to open that file in any browser with JavaScript engine and the loaded page will contain decoded value. Copy that value and paste to respond to yet request. Download will continue and complete automatically and the `decoder.html` file will be removed. 

You will need to perform that for every individual video, even if you're downloading a playlist.

#### Using Node.js to run decoding code

If an external dependency is not a problem for your use-case or you already have Node.js installed - you can progressively enhance yet with `node`. By default, yet will attempt to locate `node` binary on the system. In most cases that's sufficient and assuming you have Node.js installed - you don't need to do anything special to unlock faster download speed. 

If you'd prefer to specify `node` binary location manually, set `YET_NODE_CMD` environment variable to the full path of `node` binary (e.g. `/usr/local/bin/node` for default installation on macOS).

However you specify location of the `node` binary - upon encountering an encoded parameter, yet will download decoding code and create `decoder.js` file in the working directory and use `node` to run it and get the decoded value - this is completely automatic and doesn't require user input.

## Building for another OS

Go allows you to build binary for another OS like this (using Linux and AMD64 as an example):

```shell
env GOOS=linux GOARCH=amd64 go build -o yet
```

## Privacy

yet doesn't collect any data whatsoever. Whatever you do with yet stays on your machine. 

Your Internet connection is only used to download YouTube metadata and videos and nothing is ever uploaded anywhere. 

If you've provided any `youtube.com` cookies - they're transmitted as part of requests to get YouTube data, exactly the same way your browser would send them. You can delete `cookies.txt` that you've created at any point with no impact to ability to download publicly available videos (you won't be able to download any YouTube videos that require authorization). 
