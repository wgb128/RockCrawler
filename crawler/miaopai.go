package crawler

import (
	"github.com/xfort/RockCrawler/obj"
	"github.com/xfort/rockgo"
	"log"
	"net/url"
	"strings"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"bytes"
	"strconv"
	"regexp"
)

type MiaopaiCrawler struct {
	mphttp *rockgo.RockHttp
}

func (mp *MiaopaiCrawler) AddLog(v ...interface{}) {
	log.Println(v)
}

//读取主页所有文章数据
func (mp *MiaopaiCrawler) LoadHomeArticles(urlstr string) ([]*obj.ArticleObj, error) {

	return nil, nil
}

func (mp *MiaopaiCrawler) LoadHomeData(urlstr string, user *obj.UserObj) (suid string, userobj *obj.UserObj, err error) {
	urlObj, err := url.Parse(urlstr)
	if err != nil {
		return "", user, err
	}

	if strings.EqualFold(urlObj.Host, "www.miaopai.com") {
		urlObj.Host = "m.miaopai.com"
	}

	header := http.Header{}
	header.Set("User-Agent", rockgo.UserAgent_Android)
	resBytes, err, response := mp.mphttp.GetBytes(urlObj.String(), &header)
	if err != nil {
		return "", user, err
	}

	if response.StatusCode != 200 {
		return "", user, rockgo.NewError("读取主页网页错误")
	}

	htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(resBytes))
	if err != nil {
		return "", user, err
	}

	user.Nickname = htmlDoc.Find("b.nick").First().Text()
	user.Nickname = strings.TrimSpace(user.Nickname)
	if user.Nickname == "" {
		mp.AddLog("解析主页的nick错误", urlstr)
	}

	videoNum := htmlDoc.Find("div.box_count").First().Children().First().Text()
	videoNum = strings.TrimSpace(videoNum)
	if videoNum == "" {
		mp.AddLog("解析主页的视频数据错误", urlstr, videoNum)
	}
	user.ArticleNum, err = strconv.Atoi(videoNum)
	if err != nil {
		mp.AddLog("解析主页的视频数错误", videoNum, urlstr)
	}

	user.IconUrl, _ = htmlDoc.Find("div.head.WscaleH").First().Attr("data-url")

	scriptNodes := htmlDoc.Find("script")
	scriptLen := scriptNodes.Length()
	if scriptLen <= 0 {
		mp.AddLog("主页html内未找到script元素", urlstr, user.Nickname)
		return "", user, rockgo.NewError("主页html内未找到script元素", urlstr, user.Nickname)
	}

	for index := 0; index < scriptLen; index++ {
		itemScript, errScript := scriptNodes.Eq(index).Html()
		if errScript != nil {
			continue
		}

		itemScript = strings.TrimSpace(itemScript)
		//if strings.HasPrefix(itemScript, `FastClick.attach`) {
		user.SourceId = regexp.MustCompile("var suid = '*';").FindString(itemScript)

		//}
	}

	return user.SourceId, user, err
}