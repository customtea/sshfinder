package main

import (
	"bufio"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var regidCount int = 0

type SSHEntry struct {
	regid    int
	host     string
	hostname string
	user     string
}

const TARGET_DIR = ".ssh/"

func clearPath(path string) string {
	if len(path) > 1 && path[0:2] == "~/" {
		my, err := user.Current()
		if err != nil {
			panic(err)
		}
		path = my.HomeDir + path[1:]
	}
	path = os.ExpandEnv(path)
	return filepath.Clean(path)
}

func mapMerge(m1, m2 map[string]*SSHEntry) map[string]*SSHEntry {
	ans := map[string]*SSHEntry{}

	for k, v := range m1 {
		ans[k] = v
	}
	for k, v := range m2 {
		ans[k] = v
	}
	return (ans)
}

/*
	func check_regexp(reg, str string) bool {
		fmt.Println(regexp.MustCompile(reg).Match([]byte(str)))
		return regexp.MustCompile(reg).Match([]byte(str))
	}
*/

func loadFile(filename string) map[string]*SSHEntry {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	entries := make(map[string]*SSHEntry)
	var host string = ""

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.TrimSpace(text)
		text = strings.ToLower(text)
		if text == "" {
			continue
		} else if strings.HasPrefix(text, "#") {
			continue
		} else if strings.HasPrefix(text, "include") {
			pattern := strings.Split(text, " ")[1]
			pattern = clearPath(pattern)
			files, err := filepath.Glob(pattern)
			if err != nil {
				panic(err)
			}
			for _, file := range files {
				reslist := loadFile(file)
				entries = mapMerge(entries, reslist)
			}
		} else if strings.HasPrefix(text, "hostname") {
			hostname := strings.Split(text, " ")[1]
			entries[host].hostname = hostname
			//fmt.Println(hostname)
		} else if strings.HasPrefix(text, "host") {
			host = strings.Split(text, " ")[1]
			if strings.HasSuffix(host, "*") {
				slice := []rune(host) //文字列をスライスに変換
				host = string(slice[0 : len(slice)-1])
			}
			entries[host] = &SSHEntry{regid: regidCount, host: host, hostname: "", user: ""}
			regidCount++
			//fmt.Println(host)
		} else if strings.HasPrefix(text, "user") {
			user := strings.Split(text, " ")[1]
			entries[host].user = user
			//fmt.Println(user)
			//fmt.Println(scanner.Text())
		}
	}
	//fmt.Println(entries)
	return entries
}

func convertList(entries map[string]*SSHEntry) []*SSHEntry {
	var resList = make([]*SSHEntry, 0)

	for _, value := range entries {
		resList = append(resList, value)
	}

	sort.Slice(resList, func(i, j int) bool {
		return resList[i].regid < resList[j].regid
	})
	return resList
}

func LoadSSHConfig(listIgnore []string) []*SSHEntry {
	filename := "~/" + TARGET_DIR + "config"
	filename = clearPath(filename)
	entries := loadFile(filename)

	ignoreRegex := regexp.MustCompile("0^")
	if len(listIgnore) != 0 {
		ignorePattern := strings.Join(listIgnore, "|")
		ignoreRegex = regexp.MustCompile(ignorePattern)
	}

	//fmt.Println(ignorePattern)

	for key, value := range entries {
		if value.hostname == "" {
			delete(entries, key)
		} else if ignoreRegex.MatchString(value.host) {
			delete(entries, key)
		}
	}

	//return resDict
	resList := convertList(entries)
	return resList
}
