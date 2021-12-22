package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	listFilename  = "yet-list.txt"
	spacePrefix   = " "
	commentPrefix = "//"
)

func readList() (map[string][]string, error) {
	dirIds := make(map[string][]string)

	if _, err := os.Stat(listFilename); os.IsNotExist(err) {
		return dirIds, nil
	}

	listFile, err := os.Open(listFilename)
	if err != nil {
		return dirIds, err
	}

	scanner := bufio.NewScanner(listFile)
	scanner.Split(bufio.ScanLines)

	dir := ""

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" ||
			strings.HasPrefix(line, commentPrefix) {
			continue
		}
		if !strings.HasPrefix(line, spacePrefix) {
			dir = line
			dirIds[dir] = make([]string, 0)
		} else {
			if dir == "" {
				return dirIds, fmt.Errorf("space prefixed ids need to be under a directory (no space prefixed line)")
			}
			dirIds[dir] = append(dirIds[dir], strings.TrimSpace(line))
		}
	}

	return dirIds, nil
}
