package example

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func examplePost() {
	data := url.Values{
		"hello": {"world"},
	}

	req, err := http.NewRequest(http.MethodPost, "http://example.get.com", nil)
	if err != nil {
		log.Fatalln("err:", err)
	}
	req.PostForm = data
	resp, err := http.DefaultClient.Do(req)
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
