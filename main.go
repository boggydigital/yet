package main

import (
	"github.com/boggydigital/nod"
	"log"
	"os"
)

func main() {
	nod.EnableStdOut()

	for _, id := range os.Args[1:] {
		var err error
		if len(id) < 12 {
			err = GetVideos(id)
		} else {
			err = GetPlaylistVideos(id)
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}
