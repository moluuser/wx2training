package main

import (
	"bufio"
	"encoding/json"
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
		History [][]string `json:"history"`
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
		var roundRs []result

		for i, v := range vs {
			if strings.HasPrefix(v, me) {
				offset := 0

				var meMsgs []string
				for j, vv := range vs[i:] {
					if msg, name := getMsg(vv); name == me {
						if msg != "" {
							meMsgs = append(meMsgs, msg)
						}
						offset = j + 1
					} else {
						break
					}
				}
				// fmt.Println(meMsgs)

				var youMsgs []string
				for _, vv := range vs[i+offset:] {
					if msg, name := getMsg(vv); name == you {
						if msg != "" {
							youMsgs = append(youMsgs, msg)
						}
					} else {
						break
					}
				}
				// fmt.Println(youMsgs)

				var history [][]string
				start := 0
				if len(roundRs) > 5 {
					start = len(roundRs) / 5 * 4
				}
				for i, vv := range roundRs {
					if i < start {
						continue
					}
					history = append(history, []string{vv.Prompt, vv.Response})
				}
				if len(history) == 0 {
					history = make([][]string, 0)
				}

				if len(meMsgs) != 0 && len(youMsgs) != 0 {
					roundRs = append(roundRs, result{
						io: io{
							Prompt:   strings.Join(meMsgs, " "),
							Response: strings.Join(youMsgs, " "),
						},
						History: history,
					})
				}
				// fmt.Println(rs)
			}
		}

		rs = append(rs, roundRs...)
	}

	w, err := os.Create("result.json")
	if err != nil {
		panic(err)
	}
	for _, r := range rs {
		m, err := json.Marshal(r)
		if err != nil {
			panic(err)
		}
		w.Write(m)
		w.WriteString("\n")
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
	if n := name + ":"; strings.HasPrefix(s, n) {
		return s[len(n):], name
	}
	if index := strings.Index(s, "):"); index > 0 {
		return s[index+2:], name
	}
	return "", ""
}
