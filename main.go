package main

import (
	"encoding/json"
	"fmt"
	"github.com/lukasmwerner/pine"
	"html/template"
	"io"
	"net/http"
	"strings"
)

type ResultBundle struct {
	CommandName string `json:"commandName"`
	Parameters  string `json:"parameters"`
}

type Result struct {
	Headword   string `json:"Headword"`
	Definition string `json:"Definition"`
	ID         int    `json:"Id"`
}

type Template struct {
	templates *template.Template
}

func GetResults(words string) ([]Result, error) {
	body := strings.NewReader(fmt.Sprintf("{\"commandName\":\"cmd_autocomplete\",\"parameters\":\"%s\"}", words))
	req, err := http.NewRequest("POST", "https://www.shakespeareswords.com/ajax/AjaxResponder.aspx", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authority", "www.shakespeareswords.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Sec-Ch-Ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"98\", \"Google Chrome\";v=\"98\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://www.shakespeareswords.com")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://www.shakespeareswords.com/Public/Glossary.aspx?Id=15105")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cookie", "ASP.NET_SessionId=d3ezbb3fdecpqbbnnjgth11a; shakespeareswords.com=shwId=156b0dc4-446c-4307-b727-75967b1b1981; chkStatus=15,7")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data ResultBundle
	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}

	var results []Result

	err = json.Unmarshal([]byte(data.Parameters), &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func main() {

	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	p := pine.New()
	p.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.Query().Get("word")
		if param == "" {
			t.templates.Lookup("index.html").Execute(w, nil)
			return
		}
		results, err := GetResults(param)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		t.templates.Lookup("index.html").Execute(w, results)

	})

	http.ListenAndServe(":8080", p)
}
