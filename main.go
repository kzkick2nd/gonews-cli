package main

import (
	"flag"
	"fmt"
)

func main() {
	// OK 引数をパース
	// OK 引数でコントローラー
	// API叩く
	// 表示する
	period, command, keyword := parseArgs()
	if command == "subscribe" {
		// 登録単語を追加する
	}
	if command == "show" {
		// 登録単語と期間で検索する
	}
}

func parseArgs() (period, command, keyword string) {
	p := flag.String("past", "", "Specify time period to search for article (default=a month)")
	flag.Parse()
	return *p, flag.Arg(0), flag.Arg(1)
}
