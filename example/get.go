package example

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func exampleGet() {
	data := url.Values{}
	// or data := url.Values{"hello":[]string{"world"}}
	data.Set("hello", "world")

	parseUrl, err := url.ParseRequestURI("http://example.get.com")
	if err != nil {
		log.Fatalln("err:", err)
	}
	parseUrl.RawQuery = data.Encode()
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, parseUrl.String(), nil)
	if err != nil {
		log.Fatalln("err:", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("err:", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("err:", err)
	}
	fmt.Println(string(b))
}