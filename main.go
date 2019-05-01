package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	keywordsFilename string = ".keywords"
)

func main() {
	// OK 引数をパース
	// OK 引数でコントローラー
	// OK 登録単語リスト（読み書きOK）
	// OK 単語削除
	// OK API叩く
	// OK 表示する
	// 整形する
	// 期間指定に対応

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
	if command == "describe" {
		// 登録単語を削除する
		err := removeKeyword(listFilePath, keyword)
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

		keywords, err := readKeywords(listFilePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, v := range keywords {
			s, err := searchQuery(v)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(s)
			time.Sleep(time.Second / 2)
		}
	}
}

func parseArgs() (period, command, keyword string) {
	p := flag.String("past", "", "Specify time period to search for article. (day, week, month, default=month)")
	flag.Parse()
	return *p, flag.Arg(0), flag.Arg(1)
}

func writeKeywords(path, keyword string) error {
	w, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = fmt.Fprintf(w, "%s\n", keyword)
	if err != nil {
		return err
	}
	return nil
}

func removeKeyword(path, keyword string) error {
	var keywords string

	// Read
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() != keyword {
			keywords += scanner.Text() + "\n"
		}
	}

	// Remove
	err = os.Remove(path)
	if err != nil {
		return err
	}

	// Create
	w, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = fmt.Fprint(w, keywords)
	if err != nil {
		return err
	}
	return nil
}

func readKeywords(path string) ([]string, error) {
	var keywords []string
	f, err := os.Open(path)
	if err != nil {
		return keywords, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		keywords = append(keywords, scanner.Text())
	}

	return keywords, nil
}

// This struct formats the answer provided by the Bing News Search API.
type NewsAnswer struct {
	ReadLink     string `json: "readLink"`
	QueryContext struct {
		OriginalQuery string `json: "originalQuery"`
		AdultIntent   bool   `json: "adultIntent"`
	} `json: "queryContext"`
	TotalEstimatedMatches int `json: totalEstimatedMatches"`
	Sort                  []struct {
		Name       string `json: "name"`
		ID         string `json: "id"`
		IsSelected bool   `json: "isSelected"`
		URL        string `json: "url"`
	} `json: "sort"`
	Value []struct {
		Name  string `json: "name"`
		URL   string `json: "url"`
		Image struct {
			Thumbnail struct {
				ContentUrl string `json: "thumbnail"`
				Width      int    `json: "width"`
				Height     int    `json: "height"`
			} `json: "thumbnail"`
			Description string `json: "description"`
			Provider    []struct {
				Type string `json: "_type"`
				Name string `json: "name"`
			} `json: "provider"`
			DatePublished string `json: "datePublished"`
		} `json: "image"`
	} `json: "value"`
}

func searchQuery(term string) (string, error) {
	// Verify the endpoint URI and replace the token string with a valid subscription key.
	const endpoint = "https://japaneast.api.cognitive.microsoft.com/bing/v7.0/news/search"
	token := os.Getenv("AZURE_COGNITIVE_KEY")
	searchTerm := term

	// Declare a new GET request.
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}

	// The rest of the code in this example goes here in the main function.
	// Add the query to the request.
	param := req.URL.Query()
	param.Add("q", searchTerm)
	req.URL.RawQuery = param.Encode()

	// Insert the subscription-key header.
	req.Header.Add("Ocp-Apim-Subscription-Key", token)

	// Instantiate a client.
	client := new(http.Client)

	// Send the request to Bing.
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// Close the connection.
	defer resp.Body.Close()

	// Read the results
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Create a new answer object
	ans := new(NewsAnswer)
	err = json.Unmarshal(body, &ans)
	if err != nil {
		return "", err
	}

	// Iterate over search results and print the result name and URL.
	var text string
	for _, result := range ans.Value {
		text += result.Name + " " + result.URL + "\n"
	}
	return text, nil
}
