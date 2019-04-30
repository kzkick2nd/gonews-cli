package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	keywordsFilename string = ".keywords"
)

func main() {
	// OK 引数をパース
	// OK 引数でコントローラー
	// 登録単語リスト（読み書き）
	// API叩く
	// 表示する

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	listFilePath := filepath.Join(cwd, keywordsFilename)

	period, command, keyword := parseArgs()
	if command == "subscribe" {
		// 登録単語を追加する
		err := writeKeywords(listFilePath, keyword)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if command == "show" {
		// 登録単語と期間で検索する
		if period == "" {
			period = "month"
		}
	}
}

func parseArgs() (period, command, keyword string) {
	p := flag.String("past", "", "Specify time period to search for article. (week, month, year, default=month)")
	flag.Parse()
	return *p, flag.Arg(0), flag.Arg(1)
}

func writeKeywords(path, keyword string) error {
	w, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = fmt.Fprintf(w, " %s\n", keyword)
	return err
}
