package data

import "sync"

var ProgressMux = &sync.Mutex{}
var VideosProgress = make(map[string][]string)
