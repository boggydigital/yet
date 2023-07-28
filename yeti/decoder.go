package yeti

import (
	"bufio"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yt_urls"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	decoderStart = "function(a){var b=a.split(\"\"),"
	decoderEnd   = "return b.join(\"\")};"
)

func requestDecodedParam(hc *http.Client, n, playerPath string) (string, error) {

	// process `n` parameter:
	// 1) generate a solution file for the user
	// 2) request input from the user (they'll need to open solution file in a browser)
	// 3) return the decoded parameter to unlock fast YouTube downloads

	pu := yt_urls.PlayerUrl(playerPath)

	resp, err := hc.Get(pu.String())
	if err != nil {
		return "", err
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
		return "", err
	}

	decoderFile, err := os.Create(decoderFilename)
	if err != nil {
		return "", err
	}

	defer decoderFile.Close()

	if _, err := io.WriteString(decoderFile, "<html>"+
		"<head>"+
		"<meta name=\"color-scheme\" content=\"light dark\">"+
		"</head>"+
		"<body style='padding:1em'>"); err != nil {
		return "", err
	}
	if _, err := io.WriteString(decoderFile,
		"<output style='font-family:sans-serif;font-size:2em'>"+
			"</output>"); err != nil {
		return "", err
	}
	if _, err := io.WriteString(decoderFile, "<script>"+sb.String()+"</script>"); err != nil {
		return "", err
	}
	if _, err := io.WriteString(decoderFile,
		"<script>"+
			"document.getElementsByTagName('output')[0].textContent = "+dfn+"('"+n+"')"+
			"</script>"); err != nil {
		return "", err
	}
	if _, err := io.WriteString(decoderFile, "</body></html>"); err != nil {
		return "", err
	}

	afn, err := filepath.Abs(decoderFilename)
	if err != nil {
		return "", err
	}

	dnpa := nod.Begin("please open file://%s and paste the answer:", afn)
	defer dnpa.End()

	dn := ""
	scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		dn = scanner.Text()
		break
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if err := os.Remove(decoderFilename); err != nil {
		return "", err
	}

	return dn, nil
}
