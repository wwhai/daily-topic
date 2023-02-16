package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RollHotWord struct {
	HotWord    string `json:"hotWord"`
	SearchWord string `json:"searchWord"`
	Tag        string `json:"tag"`
	Source     string `json:"source"`
}
type TopicData struct {
	RequestId       string        `json:"requestId"`
	RollHotWordList []RollHotWord `json:"rollHotWordList"`
}

// 数据体
type Response struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    TopicData `json:"data"`
}

/*
*
* 获取网易新闻
*
 */
func get() (Response, error) {
	resp, err := http.Get("https://gw.m.163.com/search/api/v1/pc-wap/rolling-word")
	if err != nil {
		fmt.Println(err)
		return Response{}, err
	}
	defer resp.Body.Close()
	body, err0 := ioutil.ReadAll(resp.Body)
	if err0 != nil {
		fmt.Println(err0)
		return Response{}, err0
	}
	response := Response{}
	err1 := json.Unmarshal(body, &response)
	if err1 != nil {
		fmt.Println(err1)
		return Response{}, err1
	}
	return response, nil
}

/*
*
* main
*
 */
var template = `
# 每日新闻: @@DAY
## 今日热点

@@LIST

## 更多
[网易新闻] (https://www.163.com/dy/media/T1500913112740.html)
`

func main() {
	for i := 0; i < 5; i++ {
		response, err := get()
		if err == nil {
			dataList := []string{}
			// https://www.163.com/search?keyword=
			for _, element := range response.Data.RollHotWordList {
				dataList = append(dataList, "- **["+strings.Replace(element.HotWord, " ", ":", -1)+"](https://www.163.com/search?keyword="+url.QueryEscape(element.SearchWord)+")**")
			}
			today := time.Now().Format("2006-01-02 15:04:05")
			n1 := strings.Replace(template, "@@DAY", today, 1)
			n2 := strings.Replace(n1, "@@LIST", strings.Join(dataList, "\n"), 1)
			fmt.Println(n2)
			return
		}
		fmt.Println("新闻获取出错，尝试重新获取:", err)
		time.Sleep(5 * time.Second)
	}

}
