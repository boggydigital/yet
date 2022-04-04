# yet
yet is a minimalist YouTube video and channel downloader. Engineering design contraints lead to a simple application and code base. yet is built on top of `github.com/boggydigital/yt_urls`, similarly opinionated module, that provides low-level helpers to get and rationalize YouTube data.

## Using yet

```shell
yet <video-id> [, <another-video-id>...] | <playlist-id> [, <another-playlist-id>...]
```

yet supports individual video-ids or playlist-ids as args. Either one can be provided as a full `https://www.youtube.com/...` link.

Internally yet downloads videos using a list of video-ids, so any channel-ids are expanded into video-ids for download.

At the moment, there are no other (explicit) parameters that yet supports. When called without any arguments, yet can use `yet-list.txt` to update multiple playlists or if that file doesn't exist - will print help information.

### What is video-id?

Video-id is YouTube's video identifier. You can get it from a video URL : `https://www.youtube.com/watch?v=video-id`. yet can extract video-id from a URL, so you can use either video-id or `https://www.youtube.com/watch?v=video-id`.

### What is playlist-id?

Playlist-id is YouTube's videos list identifier. You can get it from a list URL: `https://www.youtube.com/watch?v=video-id&list=playlist-id`. Similarly to video-id, yet supports URL containing playlist-id, so you can use a full URL. Please note: URL can contain playlist-id and video-id at the same time and in that case playlist-id will be prioritized over video-id. If that's not desired - make sure to use either URL with just video-id or video-id itself.

## Advanced scenarios

Despite supporting only video-ids and playlist-ids, yet also (implicitly) supports channel and user videos, when they can be expressed as playlists. In addition to video-ids and playlist-ids used as arguments, yet supports `yet-list.txt` file that can be used to specify multiple sources and directories they should be downloaded. 

All advanced scenarios are detailed below. 

### Downloading channel and user videos

Most channels and user pages contain playlist links to all uploads. For example: on a channel page look for "Uploads" section and a "PLAY ALL" link - this link can be used as a playlist-id.

### Downloading videos that require authentication

Some YouTube videos require users to sign in - e.g. members-only video that user is supporting.

yet uses coost to persist session cookies and it's possible to reuse existing YouTube cookies to get access to videos that require authentication. 

Please refer to [the coost README](https://github.com/boggydigital/coost#copying-session-cookies-from-an-existing-browser-session) for the step by step guide on copying YouTube session cookies.

For yet you need to create or edit `cookies.txt` file in the yet working directory and add `youtube.com` host sessions cookies. It should look like this:

```text
youtube.com
  cookie-header=<paste-youtube-session-cookie-header-here>
```
### Specifying ffmpeg location to get the best quality video/audio

yet will attempt to locate `ffmpeg` binary on the system. If you'd prefer to specify it yourself, set `YET_FFMPEG_CMD` environment variable to the full path of `ffmpeg` binary (e.g. `/opt/homebrew/bin/ffmpeg` for Homebrew installation on macOS)
