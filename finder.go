package main

import (
	"fmt"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
)

func finder(entryList []*SSHEntry) int {
	idx, err := fuzzyfinder.Find(
		entryList,
		func(i int) string {
			return fmt.Sprintf("[%s] %s", entryList[i].host, entryList[i].hostname)
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, _ int) string {
			if i == -1 {
				return "no results"
			}
			s := fmt.Sprintf("Host %s\n    Hostname %s\n    User %s", entryList[i].host, entryList[i].hostname, entryList[i].user)
			if width < len([]rune(s)) {
				return entryList[i].host
			}
			return s
		}))
	if err == fuzzyfinder.ErrAbort {
		return -1
	}
	//fmt.Println(entryList[idx]) // The selected item.
	return idx
}
