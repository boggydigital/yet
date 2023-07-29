package yeti

import (
	"bufio"
	"fmt"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yt_urls"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	decoderStart    = "function(a){var b=a.split(\"\"),"
	decoderEnd      = "return b.join(\"\")};"
	decoderFilename = "decoder.js"
)

func decodeParam(hc *http.Client, nodeCmd, n, playerPath string) (string, error) {

	// process `n` parameter:
	// 1) generate a solution file for the user
	// 2) request input from the user (they'll need to open solution file in a browser)
	// 3) return the decoded parameter to unlock fast YouTube downloads

	dpa := nod.Begin("decoding n=%s...", n)
	defer dpa.End()

	pu := yt_urls.PlayerUrl(playerPath)

	resp, err := hc.Get(pu.String())
	if err != nil {
		return "", dpa.EndWithError(err)
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	sb := &strings.Builder{}
	ok := false
	dfn := ""

	for scanner.Scan() {
		line := scanner.Text()
		if !ok && strings.Contains(line, decoderStart) {
			if str, _, ok := strings.Cut(line, "="); ok {
				dfn = str
			}
			ok = true
		}
		if ok {
			sb.WriteString(line)
			if strings.Contains(line, decoderEnd) {
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", dpa.EndWithError(err)
	}

	decoderFile, err := os.Create(decoderFilename)
	if err != nil {
		return "", dpa.EndWithError(err)
	}

	defer decoderFile.Close()

	if _, err := io.WriteString(decoderFile, sb.String()+"\n"+
		fmt.Sprintf("console.log(%s('%s'));", dfn, n)); err != nil {
		return "", dpa.EndWithError(err)
	}

	sb.Reset()

	cmd := exec.Command(nodeCmd, decoderFilename)
	cmd.Stdout = sb
	if err := cmd.Run(); err != nil {
		return "", dpa.EndWithError(err)
	}

	dn := sb.String()

	if err := os.Remove(decoderFilename); err != nil {
		return "", dpa.EndWithError(err)
	}

	dpa.EndWithResult("done, new n=%s", dn)

	return dn, nil
}