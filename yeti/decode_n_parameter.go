package yeti

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/paths"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (

	//exceptEnhancedStr = "enhanced_except"
	maxNCodeLength = 10000
)

var (
	nDecPfx = map[string]string{
		"1": "function(a){var b=a.split(\"\"),",
		"2": "function(a){var b=String.prototype.split.call(a,\"\"),",
		"3": "function(a){var b=String.prototype.split.call(a,(\"\",\"\")),",
	}
	nDecSfx = map[string]string{
		"1": "return b.join(\"\")};",
		"2": "return Array.prototype.join.call(b,\"\")};",
		"3": "return Array.prototype.join.call(b,(\"\",\"\"))};",
	}
)

var (
	ErrNodeJsRequired        = errors.New("node.js is required")
	ErrDecoderCodeNotFound   = errors.New("decoder code not found")
	ErrDecoderCodeTooLong    = errors.New("decoder code too long")
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

	nDecoded, err := execNodeDecode(ndp, n)
	if err != nil {
		return "", dpa.EndWithError(err)
	}

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

	// buffer player code
	if _, err := buf.ReadFrom(playerContent); err != nil {
		return "", err
	}

	pfxKeys := maps.Keys(nDecPfx)
	slices.Sort(pfxKeys)
	slices.Reverse(pfxKeys)

	sfxKeys := maps.Keys(nDecSfx)
	slices.Sort(sfxKeys)
	slices.Reverse(sfxKeys)

	found := false
	var dfb, dfn string
	for _, pfxVer := range pfxKeys {
		if found {
			break
		}
		for _, sfxVer := range sfxKeys {
			dfb, dfn, err = nParamDecodeFuncBodyName(nDecPfx[pfxVer], nDecSfx[sfxVer], strings.NewReader(buf.String()))
			if err == nil && dfb != "" && dfn != "" {
				found = true
				break
			}
			if errors.Is(err, ErrDecoderCodeNotFound) || errors.Is(err, ErrDecoderCodeTooLong) {
				continue
			}
			if err != nil {
				return "", err
			}
		}
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
	okStart, okEnd := false, false
	dfn := ""

	for scanner.Scan() {
		line := scanner.Text()
		if !okStart && strings.Contains(line, start) {
			if str, _, sure := strings.Cut(line, "="); sure {
				dfn = str
			}
			okStart = true
		}
		if okStart {
			sb.WriteString(line)
			if strings.Contains(line, end) {
				okEnd = true
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", err
	}

	if sb.Len() == 0 || !okEnd {
		return "", "", ErrDecoderCodeNotFound
	}
	if sb.Len() > maxNCodeLength {
		return "", "", ErrDecoderCodeTooLong
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

func execNodeDecode(decoderPath, n string) (string, error) {
	sb := &strings.Builder{}

	cmd := exec.Command(GetBinary(NodeBin), decoderPath, n)

	cmd.Stdout = sb
	if err := cmd.Run(); err != nil {
		return "", errors.New(decoderPath + ": " + err.Error())
	}

	return strings.TrimSuffix(sb.String(), "\n"), nil
}
