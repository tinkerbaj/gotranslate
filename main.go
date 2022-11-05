package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

//It have 3 levels this is why we using 3 [][][]
type UnmarshalInterface [][][]interface{}


//From what language you want to translate -- To what language you want to translate
var (
	from = "de"
	to = "en"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	p := filepath.FromSlash(path + "/sample.txt")

	f, err := os.Open(p)
	if err != nil {
		log.Panic(err)
	}

	defer f.Close()
	scanner := bufio.NewScanner(f)
	size := 0
	linesSize := 0
	for scanner.Scan() {
		size += len([]rune(scanner.Text()))

		err := translateText(scanner.Text())
		if err != nil {
			log.Panic(err)
		}
		linesSize += 1
	}
	fmt.Printf("File contains %d characters \n and %d lines", size, linesSize)

}


func translateText(text string) error {

	//Escaped text (to make valid url request)
	enctext := url.QueryEscape(text)

	//Url without query text
	trurl := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&dt=t&sl=%s&tl=%s&q=", from, to)

	//Url with appended encryptet text
	baseurl := trurl + enctext


	//Simple GET request to Google server
	req, err := http.NewRequest(http.MethodGet, baseurl, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	//UnmarshalInterface
	var ui UnmarshalInterface
	json.Unmarshal(resBody, &ui)

	f, err := os.OpenFile("output.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	//We must cast interface to string
	if _, err := f.WriteString(ui[0][0][0].(string) + "\n"); err != nil {
		return err
	}


	return nil
}
