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

## Building for another OS

Go allows you to build binary for another OS like this (using Linux and AMD64 as an example):

```shell
env GOOS=linux GOARCH=amd64 go build -o yet
```

## Privacy

yet doesn't collect any data whatsoever. Whatever you do with yet stays on your machine. 

Your Internet connection is only used to download YouTube metadata and videos and nothing is ever uploaded anywhere. 

If you've provided any `youtube.com` cookies - they're transmitted as part of requests to get YouTube data, exactly the same way your browser would send them. You can delete `cookies.txt` that you've created at any point with no impact to ability to download publicly available videos (you won't be able to download any YouTube videos that require authorization). 
