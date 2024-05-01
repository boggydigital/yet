package yeti

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/paths"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	nParamDecoderFuncStart = "function(a){var b=a.split(\"\"),"
	nParamDecoderFuncEnd   = "return b.join(\"\")};"
)

var ErrJavaScriptRuntimeNotFound = errors.New("javascript runtime not found")

func DecodeNParam(n, playerUrl string) (string, error) {

	if !HasBinary(NodeBin) {
		return "", ErrJavaScriptRuntimeNotFound
	}

	if playerUrl == "" {
		return "", errors.New("player url is empty")
	}

	// decode `n` parameter:
	// 1) get `n` parameter decoder (local extract for a given player version, download as needed)
	// 2) run it with the JavaScript engine and capture output (transformed n)
	// 3) use the transformed parameter to unlock faster YouTube downloads

	dpa := nod.Begin("decoding n=%s...", n)
	defer dpa.End()

	ndp, err := getNParamDecoder(playerUrl)
	if err != nil {
		return "", dpa.EndWithError(err)
	}

	sb := &strings.Builder{}

	cmd := exec.Command(GetBinary(NodeBin), ndp, n)

	cmd.Stdout = sb
	if err := cmd.Run(); err != nil {
		return "", dpa.EndWithError(err)
	}

	nDecoded := strings.TrimSuffix(sb.String(), "\n")

	dpa.EndWithResult("done (n=%s)", nDecoded)

	return nDecoded, nil
}

func getNParamDecoder(playerUrl string) (string, error) {

	ndp, err := paths.AbsNParamDecoderPath(PlayerVersion(playerUrl))
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(ndp); err == nil {
		return ndp, nil
	}

	playerContent, err := GetPlayerContent(http.DefaultClient, playerUrl)
	if err != nil {
		return "", err
	}
	defer playerContent.Close()

	dfb, dfn, err := getDecodeFuncBodyName(playerContent)
	if err != nil {
		return "", err
	}

	decoderFile, err := os.Create(ndp)
	if err != nil {
		return "", err
	}
	defer decoderFile.Close()

	if err := writeNDecoder(decoderFile, dfb, dfn); err != nil {
		return "", err
	}

	return ndp, nil
}

func getDecodeFuncBodyName(r io.Reader) (string, string, error) {

	scanner := bufio.NewScanner(r)
	sb := &strings.Builder{}
	ok := false
	dfn := ""

	for scanner.Scan() {
		line := scanner.Text()
		if !ok && strings.Contains(line, nParamDecoderFuncStart) {
			if str, _, ok := strings.Cut(line, "="); ok {
				dfn = str
			}
			ok = true
		}
		if ok {
			sb.WriteString(line)
			if strings.Contains(line, nParamDecoderFuncEnd) {
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

func writeNDecoder(w io.Writer, decodeFuncBody, decodeFuncName string) error {
	if _, err := io.WriteString(w, "let "+decodeFuncBody+"\n"+
		fmt.Sprintf("console.log(%s(process.argv[process.argv.length-1]));", decodeFuncName)); err != nil {
		return err
	}
	return nil
}
