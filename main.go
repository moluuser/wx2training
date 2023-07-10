package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
	me        = ""
	you       = ""
	separator = "- - - - - - - - - - - - - - -"
)

func main() {
	type result struct {
		Instruction string `json:"instruction"`
		Input       string `json:"input"`
		Output      string `json:"output"`
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

	var (
		maxPrompt   = 0
		maxResponse = 0

		beyondLimitCount = 0
		count            = 0
		lenLimit         = 1024

		rs []result
	)
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

				prompt := strings.Join(meMsgs, " ")
				response := strings.Join(youMsgs, " ")
				if len(prompt) > maxPrompt {
					maxPrompt = len(prompt)
				}
				if len(response) > maxResponse {
					maxResponse = len(response)
				}
				if len(prompt) > lenLimit || len(response) > lenLimit {
					beyondLimitCount++
					continue
				}

				if len(meMsgs) != 0 && len(youMsgs) != 0 {
					count++
					roundRs = append(roundRs, result{
						Instruction: prompt,
						Output:      response,
					})
				}
				// fmt.Println(rs)
			}
		}

		rs = append(rs, roundRs...)
	}

	fmt.Println("maxPrompt:", maxPrompt)
	fmt.Println("maxResponse:", maxResponse)
	fmt.Println("beyondLimitCount:", beyondLimitCount)
	fmt.Println("count:", count)

	w, err := os.Create("result.json")
	if err != nil {
		panic(err)
	}
	m, err := json.Marshal(rs)
	if err != nil {
		panic(err)
	}
	w.Write(m)
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
