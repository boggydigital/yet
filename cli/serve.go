package cli

import (
	"fmt"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/rest"
	"net/http"
	"net/url"
	"strconv"
)

func ServeHandler(u *url.URL) error {
	portStr := u.Query().Get("port")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	stderr := u.Query().Has("stderr")

	return Serve(port, stderr)
}

func Serve(port int, stderr bool) error {

	if stderr {
		nod.EnableStdErrLogger()
		nod.DisableOutput(nod.StdOut)
	}

	rest.HandleFuncs()

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
