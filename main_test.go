package main

import "testing"

func TestGetMsg(t *testing.T) {
	cases := []struct {
		in   string
		msg  string
		name string
	}{
		{"[表情]", "", ""},
		{"[图片]", "", ""},
		{"Mångata🌙:我给你买了个好东西", "我给你买了个好东西", you},
		{"Mångata🌙 (2022-04-11 12:04:55):我也不知道", "我也不知道", you},
		{"", "", ""},
		{"- - - - - - - - - - - - - - -", "", ""},
		{"不太好搞", "", ""},
		{"Lollipop (2022-04-11 13:42:25):地上看看？", "地上看看？", me},
		{"Lollipop:晚上暂时打算屯点吃的东西", "晚上暂时打算屯点吃的东西", me},
	}

	for _, c := range cases {
		msg, name := getMsg(c.in)
		if msg != c.msg || name != c.name {
			t.Errorf("getMsg(%q) == %q, %q, want %q, %q", c.in, msg, name, c.msg, c.name)
		}
	}
}
