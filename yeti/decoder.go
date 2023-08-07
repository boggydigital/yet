package yeti

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yt_urls"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	decoderStart           = "function(a){var b=a.split(\"\"),"
	decoderEnd             = "return b.join(\"\")};"
	decoderNodeFilename    = "decoder.js"
	decoderBrowserFilename = "decoder.html"
)

var memoizer = make(map[string]string)

func decodeParam(hc *http.Client, nodeCmd, n, playerUrl string) (string, error) {

	if dn, ok := memoizer[n+playerUrl]; ok {
		return dn, nil
	}

	if playerUrl == "" {
		return "", errors.New("player url is empty")
	}

	// transform `n` parameter:
	// 1) generate a file containing player specific transform function
	// 2) run it with the Node.js and capture output (transformed n)
	// 3) use the transformed parameter to unlock faster YouTube downloads

	dpa := nod.Begin("decoding n=%s...", n)
	defer dpa.End()

	pu := yt_urls.PlayerUrl(playerUrl)

	resp, err := hc.Get(pu.String())
	if err != nil {
		return "", dpa.EndWithError(err)
	}

	defer resp.Body.Close()

	dfb, dfn, err := getDecodeFuncBodyName(resp.Body)
	if err != nil {
		return "", dpa.EndWithError(err)
	}

	filename := decoderNodeFilename
	if nodeCmd == "" {
		filename = decoderBrowserFilename
	}

	decoderFile, err := os.Create(filename)
	if err != nil {
		return "", dpa.EndWithError(err)
	}

	defer decoderFile.Close()

	var createFile func(io.Writer, string, string, string) error
	switch filename {
	case decoderNodeFilename:
		createFile = createNodeFile
	case decoderBrowserFilename:
		createFile = createBrowserFile
	}

	if err := createFile(decoderFile, dfb, dfn, n); err != nil {
		return "", dpa.EndWithError(err)
	}

	var decoder func(string, string) (string, error)
	if nodeCmd != "" {
		decoder = decodeWithNode
	} else {
		decoder = decodeWithBrowser
	}

	dn, err := decoder(filename, nodeCmd)
	if err != nil {
		return "", dpa.EndWithError(err)
	}

	dn = strings.TrimSuffix(dn, "\n")

	memoizer[n+playerUrl] = dn

	if err := os.Remove(filename); err != nil {
		return "", dpa.EndWithError(err)
	}

	dpa.EndWithResult("done (n=%s)", dn)

	return dn, nil
}

func getDecodeFuncBodyName(r io.Reader) (string, string, error) {

	scanner := bufio.NewScanner(r)
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
		return "", "", err
	}

	if sb.Len() == 0 {
		return "", "", errors.New("decoder code not found")
	}

	return sb.String(), dfn, nil
}

func createNodeFile(w io.Writer, decodeFuncBody, decodeFuncName, n string) error {
	if _, err := io.WriteString(w, decodeFuncBody+"\n"+
		fmt.Sprintf("console.log(%s('%s'));", decodeFuncName, n)); err != nil {
		return err
	}
	return nil
}

func createBrowserFile(w io.Writer, decodeFuncBody, decodeFuncName, n string) error {
	if _, err := io.WriteString(w, "<html>"+
		"<head>"+
		"<meta name=\"color-scheme\" content=\"light dark\">"+
		"</head>"+
		"<body style='padding:1em'>"); err != nil {
		return err
	}
	if _, err := io.WriteString(w,
		"<output style='font-family:sans-serif;font-size:2em'>"+
			"</output>"); err != nil {
		return err
	}
	if _, err := io.WriteString(w, fmt.Sprintf("<script>%s</script>", decodeFuncBody)); err != nil {
		return err
	}
	if _, err := io.WriteString(w,
		"<script>"+
			fmt.Sprintf("document.getElementsByTagName('output')[0].textContent = %s('%s')", decodeFuncName, n)+
			"</script>"); err != nil {
		return err
	}
	if _, err := io.WriteString(w, "</body></html>"); err != nil {
		return err
	}
	return nil
}

func decodeWithNode(filename, nodeCmd string) (string, error) {

	sb := &strings.Builder{}

	cmd := exec.Command(nodeCmd, filename)
	cmd.Stdout = sb
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

func decodeWithBrowser(filename, _ string) (string, error) {

	afn, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}

	dnpa := nod.Begin("please open file://%s and paste the answer:", afn)
	defer dnpa.End()

	dn := ""

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		dn = scanner.Text()
		break
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return dn, nil
}
