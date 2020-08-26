package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	// -hオプション用文言
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s [options...] <url>\n", os.Args[0])
		flag.PrintDefaults()
	}
	
	var requestHeader bool
	flag.BoolVar(&requestHeader, "v", false, "-v, --verbose  Make the operation more talkative")
	var outputFile string
	flag.StringVar(&outputFile, "o", "", "-o, --output <file>  Write to file instead of stdout")
	var postValues string
	flag.StringVar(&postValues, "d", "", "-d, --data <data>  HTTP POST data")
	var requestType string
	flag.StringVar(&requestType, "X", "GET", "-X, --request <command>  Specify request command to use")
	flag.Parse()



	// URLの指定がない時
	if len(flag.Args())<=0{
		fmt.Printf("%s: no URL specified!\n", os.Args[0])
		fmt.Printf("%s: try '%s --help' or '%s --manual' for more information\n", os.Args[0], os.Args[0], os.Args[0])
		os.Exit(1)
	}

	if contains(os.Args, "-d") && postValues==""{
		fmt.Printf("%s: option -d: requires parameter\n", os.Args[0])
		fmt.Printf("%s: try '%s --help' or '%s --manual' for more information\n", os.Args[0], os.Args[0], os.Args[0])
		os.Exit(1)
	}

	addr := flag.Arg(0)
	// postValuesをsplitしてurl.Values{}に格納
	values := url.Values{}
	if postValues != ""{
		splitEach := strings.Split(postValues, "&")
		for _, v := range splitEach{
			splitKeyVaue := strings.Split(v, "=")
			values.Add(splitKeyVaue[0], splitKeyVaue[1])
		}
		fmt.Println(values.Encode())
	}

	if requestType == "GET"{
		get(addr, requestHeader, outputFile)
	}else if requestType == "POST"{
		post(addr, requestHeader, values)
	}else{
		fmt.Printf("%s: requestType is not correct!\n", os.Args[0])
		fmt.Printf("%s: try 'kcurl --help' or 'kcurl --manual' for more information\n", os.Args[0])
		os.Exit(1)
	}

}

func get(addr string, requestHeader bool, filename string){
	// レスポンスを作成
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: tr,
	}
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	req.Header.Add("If-None-Match", `W/"wyzzy"`)

	// リクエストを送信
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// レスポンスを受信して表示
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// TODO: ここのエラー処理直したい
	// -oオプションしかなくファイル名が指定されていない時はエラー表示
	if contains(os.Args, "-o") && filename == ""{
		fmt.Printf("%s: option -o: requires parameter\n", os.Args[0])
		fmt.Printf("%s: try '%s --help' or '%s --manual' for more information\n", os.Args[0], os.Args[0], os.Args[0])
		os.Exit(1)
	}

	// -oオプションがあってファイル名が指定されているときのみファイル作成
	if contains(os.Args, "-o") && filename != ""{
		fp, err := os.Create(filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer fp.Close()

		fp.WriteString(string(responseBody))
	}

	// -vオプションがあるときリクエスト内容を表示
	if requestHeader{
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dump))
	}

	fmt.Println(resp.Status)
	fmt.Println(string(responseBody))
}

func post(addr string, requestHeader bool, values url.Values){


	// レスポンスを作成
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: tr,
	}
	req, err := http.NewRequest("POST", addr, strings.NewReader(values.Encode()))
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// -vオプションがあるときリクエスト内容を表示
	if requestHeader{
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dump))
	}

	// リクエストを送信
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// レスポンスを受信して表示
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Status)
	fmt.Println(string(responseBody))
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}
