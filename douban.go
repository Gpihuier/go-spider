package go_spider

import (
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
)

const (
	URL       = "https://movie.douban.com/top250"
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"
)

type Result struct {
	Index       string `json:"index"`        // 序号
	Picture     string `json:"picture"`      // 图片
	Link        string `json:"link"`         // 链接
	Title       string `json:"title"`        // 标题
	Playable    bool   `json:"playable"`     // 是否可播放
	Description string `json:"description"`  // 描述
	Appraise    string `json:"appraise"`     // 评分
	AppraiseNum string `json:"appraise_num"` // 评价人数
	HotAppraise string `json:"hot_appraise"` // 亮评
}

type Spider struct {
	Url       string
	UserAgent string
}

func NewSpider() *Spider {
	return &Spider{
		Url:       URL,
		UserAgent: UserAgent,
	}
}

// Run 运行爬虫
func (s *Spider) Run() ([]Result, error) {
	res := make([]Result, 25)
	req, err := http.NewRequest(http.MethodGet, s.Url, nil)
	if err != nil {
		return res, err
	}
	req.Header.Set("user-agent", s.UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	// each html
	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return res, err
	}

	nodes := htmlquery.Find(doc, "//div[@class='article']//ol/li")
	for k, node := range nodes {
		res[k].Index = htmlquery.FindOne(node, "./div[@class='item']/div[@class='pic']/em/text()").Data
		res[k].Picture = htmlquery.InnerText(htmlquery.FindOne(node, "./div[@class='item']/div[@class='pic']/a/img/@src"))
		res[k].Link = htmlquery.InnerText(htmlquery.FindOne(node, "./div[@class='item']/div[@class='info']/div[@class='hd']/a/@href"))
		res[k].Title = htmlquery.InnerText(htmlquery.FindOne(node, "./div[@class='item']/div[@class='info']/div[@class='hd']/a/span/text()"))
		isPlay := htmlquery.FindOne(node, "./div[@class='item']/div[@class='info']/div[@class='hd']/span[@class='playable']/text()")
		if isPlay != nil {
			res[k].Playable = true
		}
		description := htmlquery.Find(node, "./div[@class='item']/div[@class='info']/div[@class='bd']/p[1]/text()")
		for _, desc := range description {
			res[k].Description += strings.TrimSpace(htmlquery.InnerText(desc)) + ","
		}
		res[k].Description = strings.TrimSuffix(res[k].Description, ",")
		res[k].Appraise = htmlquery.InnerText(
			htmlquery.FindOne(node, "./div[@class='item']/div[@class='info']/div[@class='bd']/div[@class='star']/span[@class='rating_num']/text()"))
		res[k].AppraiseNum = htmlquery.InnerText(htmlquery.FindOne(node, "./div[@class='item']/div[@class='info']/div[@class='bd']/div[@class='star']/span[last()]/text()"))
		hasHotAppraise := htmlquery.FindOne(node, "./div[@class='item']/div[@class='info']/div[@class='bd']/p[last()]/span/text()")
		if hasHotAppraise != nil {
			res[k].AppraiseNum = htmlquery.InnerText(hasHotAppraise)
		} else {
			res[k].AppraiseNum = "暂无亮评"
		}
	}
	return res, err
}
