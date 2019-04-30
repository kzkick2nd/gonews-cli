package main

import (
	"flag"
	"fmt"
)

func main() {
	// 引数をパース
	// API叩く
	// 表示する
	period, command, keyword := parseArgs()
	fmt.Println(period, command, keyword)
}

func parseArgs() (period, command, keyword string) {
	p := flag.String("past", "", "Specify time period to search for article (default=a month)")
	flag.Parse()
	return *p, flag.Arg(0), flag.Arg(1)
}
