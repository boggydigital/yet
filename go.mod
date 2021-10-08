module yet

go 1.17

require (
	github.com/boggydigital/dolo v0.1.4-alpha
	github.com/boggydigital/nod v0.1.4
	github.com/boggydigital/yt_urls v0.1.1
)

require (
	github.com/boggydigital/match_node v0.1.1 // indirect
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
)

replace (
	github.com/boggydigital/dolo => ../dolo
	github.com/boggydigital/nod => ../nod
	github.com/boggydigital/yt_urls => ../yt_urls
)
