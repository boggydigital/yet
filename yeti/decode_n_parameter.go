package yeti

import (
	"bufio"
	"bytes"
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

var (
	nDecPfx = map[string]string{
		"v1": "function(a){var b=a.split(\"\"),",
		"v2": "function(a){var b=String.prototype.split.call(a,\"\"),",
	}
	nDecSfx = map[string]string{
		"v1": "return b.join(\"\")};",
		"v2": "return Array.prototype.join.call(b,\"\")};",
	}
)

var (
	ErrNodeJsRequired        = errors.New("node.js is required")
	ErrDecoderCodeNotFound   = errors.New("decoder code not found")
	ErrNParamDecoderNotFound = errors.New("n-param decoder not found")
)

func DecodeNParam(n, playerUrl string) (string, error) {

	if n == "" {
		return "", nil
	}

	if !HasBinary(NodeBin) {
		return "", ErrNodeJsRequired
	}

	if playerUrl == "" {
		return "", errors.New("player url is empty")
	}

	// decode `n` parameter:
	// 1) get `n` parameter decoder (local extract for a given player version, download as needed)
	// 2) run it with the Node.js and capture output (transformed n)
	// 3) use the transformed parameter to unlock faster YouTube downloads

	dpa := nod.Begin("decoding n=%s...", n)
	defer dpa.End()

	ndp, err := getNParamDecoder(playerUrl)
	if err != nil {
		return "", dpa.EndWithError(err)
	}

	nDecoded, err := execNodeDecodeNParam(ndp, n)

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

	bts := make([]byte, 0)
	buf := bytes.NewBuffer(bts)
	tr := io.TeeReader(playerContent, buf)

	dfb, dfn, err := nParamDecodeFuncBodyName(nDecPfx["v1"], nDecSfx["v1"], tr)
	if errors.Is(err, ErrDecoderCodeNotFound) {
		dfb, dfn, err = nParamDecodeFuncBodyName(nDecPfx["v2"], nDecSfx["v2"], buf)
	}

	if err != nil {
		return "", err
	}

	if dfb == "" || dfn == "" {
		return "", ErrNParamDecoderNotFound
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

func nParamDecodeFuncBodyName(start, end string, r io.Reader) (string, string, error) {

	scanner := bufio.NewScanner(r)
	sb := &strings.Builder{}
	ok := false
	dfn := ""

	for scanner.Scan() {
		line := scanner.Text()
		if !ok && strings.Contains(line, start) {
			if str, _, ok := strings.Cut(line, "="); ok {
				dfn = str
			}
			ok = true
		}
		if ok {
			sb.WriteString(line)
			if strings.Contains(line, end) {
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", err
	}

	if sb.Len() == 0 {
		return "", "", ErrDecoderCodeNotFound
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

func execNodeDecodeNParam(decoderPath, n string) (string, error) {
	sb := &strings.Builder{}

	cmd := exec.Command(GetBinary(NodeBin), decoderPath, n)

	cmd.Stdout = sb
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSuffix(sb.String(), "\n"), nil
}
