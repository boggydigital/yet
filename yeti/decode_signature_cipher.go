package yeti

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/paths"
	"github.com/boggydigital/yt_urls"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	signatureCipherFuncBodyStart = "a=a.split(\"\")"
	objectStartTemplate          = "var %s={"
	objectEnd                    = "};"
)

var (
	ErrSignatureCipherFunctionNotFound = errors.New("signature cipher function not found")
	ErrSignatureCipherObjectNotFound   = errors.New("signature cipher object not found")
)

func DecodeSignatureCiphers(hc *http.Client, ipr *yt_urls.InitialPlayerResponse) error {
	if !ipr.SignatureCipher() {
		return nil
	}

	if !HasBinary(NodeBin) {
		return ErrNodeJsRequired
	}

	// decode signatureCiphers:
	// 1) get `signatureCipher` decoder (local extract for a given player version, download as needed)
	// 2) for each URL containing signatureCipher:
	// 3) run decoder with the Node.js and capture output (decoded signatureCipher)
	// 3) add .Url param to formats and continue to use InitialPlayerResponse as usual

	dsca := nod.Begin("decoding signatureCiphers...")
	defer dsca.End()

	scdp, err := getSignatureCipherDecoder(ipr.PlayerUrl)
	if err != nil {
		return dsca.EndWithError(err)
	}

	// decoding three formats most likely to be used by yet
	formats := []*yt_urls.Format{
		ipr.BestFormat(), ipr.BestAdaptiveVideoFormat(), ipr.BestAdaptiveAudioFormat(),
	}

	for i, f := range formats {
		u, err := decodeSignatureCipher(f.SignatureCipher, scdp)
		if err != nil {
			return dsca.EndWithError(err)
		}
		formats[i].Url = u.String()
	}

	return nil
}

func decodeSignatureCipher(signatureCipher, decoderPath string) (*url.URL, error) {
	scq, err := url.ParseQuery(signatureCipher)
	if err != nil {
		return nil, err
	}

	decodedSignature, err := execNodeDecodeNParam(decoderPath, scq.Get("s"))
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(scq.Get("url"))
	if err != nil {
		return nil, err
	}

	nq := u.Query()
	nq.Add("sig", decodedSignature)

	u.RawQuery = nq.Encode()

	return u, nil
}

func getSignatureCipherDecoder(playerUrl string) (string, error) {
	scdp, err := paths.AbsSignatureCipherDecoderPath(PlayerVersion(playerUrl))
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(scdp); err == nil {
		return scdp, nil
	}

	playerContent, err := GetPlayerContent(http.DefaultClient, playerUrl)
	if err != nil {
		return "", err
	}
	defer playerContent.Close()

	scfb, scfn, err := signatureCipherFuncBodyName(playerContent)
	if err != nil {
		return "", err
	}

	if scfb == "" || scfn == "" {
		return "", ErrSignatureCipherFunctionNotFound
	}

	scob, err := signatureCipherObjectBody(playerContent, signatureCipherObjectName(scfb))

	if scob == "" {
		return "", ErrSignatureCipherObjectNotFound
	}

	decoderFile, err := os.Create(scdp)
	if err != nil {
		return "", err
	}
	defer decoderFile.Close()

	if err := writeSignatureCipherDecoder(decoderFile, scfn, scfb, scob); err != nil {
		return "", err
	}

	return scdp, nil
}

func signatureCipherFuncBodyName(r io.Reader) (string, string, error) {

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, signatureCipherFuncBodyStart) {
			if scfn, _, ok := strings.Cut(line, "="); ok {
				return line, scfn, nil
			}
		}
		if err := scanner.Err(); err != nil {
			return "", "", err
		}
	}

	return "", "", ErrSignatureCipherFunctionNotFound
}

func signatureCipherObjectName(funcBody string) string {

	if parts := strings.Split(funcBody, ";"); len(parts) > 1 {
		if on, _, ok := strings.Cut(parts[1], "."); ok {
			return on
		}
	}
	return ""
}

func signatureCipherObjectBody(r io.Reader, objectName string) (string, error) {
	scanner := bufio.NewScanner(r)
	sb := &strings.Builder{}
	ok := false
	objectStart := fmt.Sprintf(objectStartTemplate, objectName)

	for scanner.Scan() {
		line := scanner.Text()

		if ok {
			if !strings.Contains(line, objectEnd) {
				sb.WriteString(line)
			} else {
				if remainder, _, ok := strings.Cut(line, objectEnd); ok {
					sb.WriteString(remainder)
					sb.WriteString(objectEnd)
				}
				break
			}
		}

		if !ok && strings.Contains(line, objectStart) {
			sb.WriteString(objectStart)
			if _, trailer, ok := strings.Cut(line, objectStart); ok {
				sb.WriteString(trailer)
			}
			ok = true
		}

		if err := scanner.Err(); err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}

func writeSignatureCipherDecoder(w io.Writer, funcName, funcBody, objectBody string) error {
	if _, err := io.WriteString(w, objectBody+"\n"+
		"let "+funcBody+"\n"+
		fmt.Sprintf("console.log(%s(process.argv[process.argv.length-1]));", funcName)); err != nil {
		return err
	}
	return nil
}
