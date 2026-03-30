package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"math/rand"
	"moxin-zhicheng/internal/config"
	"moxin-zhicheng/internal/database"
	"moxin-zhicheng/internal/logger"
	model "moxin-zhicheng/internal/models"
	"net/url"
	"regexp"
	"strings"
	"time"
)

//type PoetryData struct {
//	Title       string
//	Author      string
//	Translation string
//	Annotation  string
//	Rhythmic    string
//}

func main() {
	// 1. 初始化数据库连接
	logger.InitLogger("dev")
	config.InitConfig()
	database.InitDB()

	// 2. 获取待处理的数据（译文或注释为空的）
	var poetries []model.Poetry
	database.DB.Where("translation is null OR annotation is null").Find(&poetries)

	for _, p := range poetries {
		// 优先使用词牌名，没有则用标题
		if p.Type == "caocao" {
			p.Author = "曹操"
			fmt.Printf("[修正] 检测到 type 为 caocao，已将作者修正为: %s\n", p.Author)
		}
		searchKey := p.Rhythmic
		if searchKey == "" {
			searchKey = p.Title
		}

		fmt.Printf("\n[任务] 正在处理: 《%s》 - %s\n", searchKey, p.Author)

		// 3. 调用爬虫获取原始 HTML
		data, err := ScrapeGushiwen(searchKey, p.Author)
		if err != nil || (data.Translation == "" && data.Annotation == "") {
			fmt.Printf("[跳过] 未找到匹配内容或请求失败: %v\n", err)
			continue
		}

		// 4. 清洗数据并更新
		p.Translation = cleanHTML(data.Translation)
		p.Annotation = cleanHTML(data.Annotation)

		// 保存回数据库
		if err := database.DB.Save(&p).Error; err != nil {
			fmt.Printf("[错误] 数据库更新失败: %v\n", err)
		} else {
			fmt.Println("[成功] 已更新译文与注释")
		}

		// 适当延时，保护 IP 不被封
		time.Sleep(time.Duration(4+rand.Intn(3)) * time.Second)
	}
}

// cleanHTML 专门负责把爬到的 HTML 片段转为干净的纯文本
func cleanHTML(htmlContent string) string {
	if htmlContent == "" {
		return ""
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return ""
	}
	// 针对古诗文网的结构提取正文
	return strings.TrimSpace(doc.Find("div.contson").Text())
}

func ScrapeGushiwen(title, author string) (*model.Poetry, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("so.gushiwen.cn", "www.gushiwen.cn"),
	)
	extensions.RandomUserAgent(c)
	result := &model.Poetry{Title: title, Author: author}

	// 统一设置 Header
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "login=flase; ticketStr=205109234%7cgQGU8TwAAAAAAAAAAS5odHRwOi8vd2VpeGluLnFxLmNvbS9xLzAyTi1rQ1FfbGVkN2kxT3B2TGhGMUcAAgSZksdpAwQAjScA; wxopenid=oVc5H0sNPp-AR7TsWBdvnw2QwWjw; Hm_lvt_9007fab6814e892d3020a64454da5a55=1773833062,1774157161,1774686861,1774694927; HMACCOUNT=523E7D76752CFD9F; Hm_lpvt_9007fab6814e892d3020a64454da5a55=1774695082; gsw2017user=7795777%7c6A5471B38CFFFF27880E4F7E9679CF7A537e2a7a%7c2000%2f1%2f1%7c2000%2f1%2f1; userPlay=7795777%7C0%7C0%7C2%7C1%7C0%7C0%7C1%7C0%7C0%7C0%7C0%7C0%7C0")
		//r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

		if strings.Contains(r.URL.String(), "ajax") {
			r.Headers.Set("Referer", "https://www.gushiwen.cn/shiwenv_1921957e6e83.aspx")
		} else {
			r.Headers.Set("Referer", "https://www.gushiwen.cn/")
		}
	})

	// 处理搜索页，提取加密ID并请求注释和译文
	c.OnHTML("div.zongheShiwen", func(e *colly.HTMLElement) {
		e.ForEach("img[id^='btnYiwen']", func(_ int, el *colly.HTMLElement) {
			onclick := el.Attr("onclick")

			// 正则匹配第二个参数，即 32 位的加密 ID
			re := regexp.MustCompile(`OnYiwenSearch\('[^']+',\s*'([^']+)'`)
			match := re.FindStringSubmatch(onclick)

			if len(match) > 1 {
				targetID := match[1]
				fmt.Printf("[DEBUG] 捕获到加密 ID: %s\n", targetID)

				// 请求注释（value=zhu）
				if result.Annotation == "" {
					annotationURL := fmt.Sprintf("https://www.gushiwen.cn/nocdn/ajaxshiwenDetailCont.aspx?id=%s&value=zhu", targetID)
					fmt.Println("[DEBUG] 正在请求注释...")
					err := e.Request.Visit(annotationURL)
					if err != nil {
						fmt.Printf("[ERROR] 注释请求失败: %v\n", err)
					}
				}

				// 请求译文（value=yi）
				if result.Translation == "" {
					translationURL := fmt.Sprintf("https://www.gushiwen.cn/nocdn/ajaxshiwenDetailCont.aspx?id=%s&value=yi", targetID)
					fmt.Println("[DEBUG] 正在请求译文...")
					err := e.Request.Visit(translationURL)
					if err != nil {
						fmt.Printf("[ERROR] 译文请求失败: %v\n", err)
					}
				}
			}
		})
	})

	// 统一处理接口响应（注释和译文）
	c.OnResponse(func(r *colly.Response) {
		urlStr := r.Request.URL.String()

		if strings.Contains(urlStr, "ajaxshiwenDetailCont.aspx") {
			fmt.Printf("[NETWORK] 收到接口响应 | 状态码: %d | 数据长度: %d 字节\n", r.StatusCode, len(r.Body))

			if r.StatusCode != 200 {
				fmt.Printf("[ERROR] 接口请求非 200，实际状态: %d\n", r.StatusCode)
				return
			}

			body := strings.TrimSpace(string(r.Body))

			if len(body) < 50 {
				if len(body) == 0 {
					fmt.Println("[WARN] 接口返回内容完全为空！可能原因：Cookie 域名不匹配、Referer 缺失、或被频率限制")
				} else {
					fmt.Printf("[WARN] 接口响应异常短。收到原始内容: [%s]\n", body)
				}
				fmt.Printf("[DEBUG] 当前请求 Headers: %v\n", r.Request.Headers)
				return
			}

			// 根据 value 参数判断是注释还是译文
			if strings.Contains(urlStr, "value=zhu") {
				result.Annotation = body
				fmt.Println("[DEBUG] ✨ 注释解析成功，已存入 PoetryData")
			} else if strings.Contains(urlStr, "value=yi") {
				result.Translation = body
				fmt.Println("[DEBUG] ✨ 译文解析成功，已存入 PoetryData")
			}
		}
	})

	// 监听网络层面的错误
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("[FATAL NETWORK ERROR] 访问 %s 失败: %v\n", r.Request.URL, err)
	})

	// 启动搜索
	searchURL := fmt.Sprintf("https://www.gushiwen.cn/search.aspx?value=%s+%s",
		url.QueryEscape(title), url.QueryEscape(author))

	fmt.Printf("[DEBUG] 启动初始搜索: %s\n", searchURL)
	err := c.Visit(searchURL)

	if err != nil {
		fmt.Printf("[ERROR] 初始 Visit 失败: %v\n", err)
	}

	return result, err
}
