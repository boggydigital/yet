# yet
yet is a Minimalist YouTube video and channel downloader. Engineering design contraints lead to a simple application and code base. 

##Using yet

```shell
yet [video-id, ...] | [playlist-id...]
```

yet supports individual video-ids or playlist-ids as args. Either one can be provided as a full `https://www.youtube.com/...` link.

Internally yet downloads videos using a list of video-ids, so any channel-ids are expanded into video-ids for download.

At the moment, there are no other parameters that yet supports.

###What is video-id?

Video-id is YouTube's video identifier. You can get it from a video UR : `https://www.youtube.com/watch?v=video-id`

###What is playlist-id?

Playlist-id is YouTube's videos list identifier. You can get it from a list URL: `https://www.youtube.com/watch?v=video-id&list=playlist-id`

##Advanced Scenarios

Despite supporting only video-ids and playlist-ids, yet also (implicitly) supports channel and user videos, when they can be expressed as playlists.

### Downloading channel and user videos

Most channels and user pages contain playlist links to all uploads. For example: on a channel page look for "Uploads" section and a "▶ PLAY ALL" link - this link can be used as a playlist-id.

### Downloading videos that require authentication

TBD

##Backlog

TBD