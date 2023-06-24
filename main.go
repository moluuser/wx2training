package main

import (
	"bufio"
	"os"
	"strings"
)

type void struct{}

var (
	ignoredWords = map[string]void{
		"[表情]":      void{},
		"[图片]":      void{},
		"[视频/语音通话]": void{},
		"[语音":       void{},
		"[文件":       void{},
	}
)

const (
	me        = "Lollipop"
	you       = "Mångata🌙"
	separator = "- - - - - - - - - - - - - - -"
)

func main() {
	type io struct {
		Prompt   string `json:"prompt"`
		Response string `json:"response"`
	}

	type result struct {
		io
		History []io `json:"history"`
	}

	var records []string
	r, err := os.Open("raw.txt")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	s := bufio.NewScanner(r)
	for s.Scan() {
		records = append(records, s.Text())
	}

	var rounds [][]string
	rounds = append(rounds, []string{})
	for _, v := range records {
		if v == separator {
			rounds = append(rounds, []string{})
			continue
		}
		rounds[len(rounds)-1] = append(rounds[len(rounds)-1], v)
	}

	var rs []result
	for _, vs := range rounds {
		for i, v := range vs {
			if strings.HasPrefix(v, me) {
				offset := 0

				var meMsgs []string
				for j, vv := range vs[i:] {
					if msg, name := getMsg(vv); name == me {
						meMsgs = append(meMsgs, msg)
						offset = j + 1
					} else {
						break
					}
				}
				// fmt.Println(meMsgs)

				var youMsgs []string
				for _, vv := range vs[i+offset:] {
					if msg, name := getMsg(vv); name == you {
						youMsgs = append(youMsgs, msg)
					} else {
						break
					}
				}
				// fmt.Println(youMsgs)

				if len(meMsgs) != 0 && len(youMsgs) != 0 {
					rs = append(rs, result{
						io: io{
							Prompt:   strings.Join(meMsgs, " "),
							Response: strings.Join(youMsgs, " "),
						},
					})
				}
				// fmt.Println(rs)
			}
		}
	}
}

func getMsg(s string) (msg string, name string) {
	isFind := false
	for k := range ignoredWords {
		if strings.Contains(s, k) {
			isFind = true
			break
		}
	}
	if isFind {
		return "", ""
	}

	isMe := false
	isYou := false
	if strings.HasPrefix(s, me) {
		isMe = true
	}
	if strings.HasPrefix(s, you) {
		isYou = true
	}
	if isMe {
		name = me
	} else if isYou {
		name = you
	}
	if name == "" {
		return "", ""
	}
	if n := name + ":"; strings.HasPrefix(s, n) {
		return s[len(n):], name
	}
	if index := strings.Index(s, "):"); index > 0 {
		return s[index+2:], name
	}
	return "", ""
}
