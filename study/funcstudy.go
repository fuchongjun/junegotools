package study

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/goquery"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"path"
	"strings"
	"time"
)

var htmlDocLinks []string

//大数字加逗号
func comma(s string) string {
	var buf bytes.Buffer
	slen := len(s)
	for i := 1; i <= slen; i++ {
		buf.WriteByte(s[i-1])

		if ((slen-i)%3 == 0) && (i != slen) {
			buf.WriteByte(',')
		}
	}
	return buf.String()

}

//字符串同文异构
func YiGou(str1, str2 string) bool {
	if len(str1) != len(str2) {
		return false
	}
	var map1 = make(map[string]int, len(str1))

	for _, v := range str1 {
		map1[string(v)] += 1
	}
	var map2 = make(map[string]int, len(str2))

	for _, v := range str2 {
		map2[string(v)] += 1
	}
	for k, v := range map1 {
		if map2[k] != v {
			return false
		}
	}
	return true

}

//遍历html文档树节点
func Visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}

		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = Visit(links, c)
	}
	return links
}
func VisitImg(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "img" {
		for _, a := range n.Attr {
			if a.Key == "src" {
				links = append(links, a.Val)
			}

		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = VisitImg(links, c)
	}
	return links

}

//函数作为参数，重写visit--学习递归
func ForEachNode(n *html.Node, pre, post func(node *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ForEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
func getNodeLinks(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, item := range node.Attr {
			if item.Key == "href" {
				htmlDocLinks = append(htmlDocLinks, item.Val)
			}
		}
	}
}

//按空格分割获取到html内容
func SpliteHtmByWord() {
	var url = "https://www.cnblogs.com/"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reading url error:%v\n", err)
		os.Exit(1)
	}
	scaner := bufio.NewScanner(resp.Body)
	scaner.Split(bufio.ScanWords)
	count := 0
	for scaner.Scan() {
		fmt.Println(scaner.Text())
		count++
	}
	if err := scaner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
	resp.Body.Close()
	fmt.Printf("%d\n", count)
}

//获取html文档中的所有连接
func GetHtmlDocLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		err = fmt.Errorf("网络请求错误，请求URL:%s,错误信息：%v", url, err)
		return nil, err
	}
	nodes, err := html.Parse(resp.Body)
	if err != nil {
		err = fmt.Errorf("文档解析发生错误：%v", err)
		return nil, err
	}
	ForEachNode(nodes, getNodeLinks, nil)
	if htmlDocLinks == nil {
		err = fmt.Errorf("遍历html文档节点发生错误，或无a标签：%v", err)
		return nil, err
	}
	return htmlDocLinks, err
}

//过滤html文档中单词个数和图片个数
func CountWordsAndImages() (wordcount, imagecount int) {
	const url = "http://www.honliv.com.cn/"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Errorf("请求发生错误：%s", err)
		return
	}

	htmsbytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	htmdata := make([]byte, len(htmsbytes))
	copy(htmdata, htmsbytes)
	htmlNodes, err := html.Parse(bytes.NewReader(htmsbytes))
	if err != nil {
		fmt.Errorf("html转换发生错误：%s", err)
		return
	}
	var images []string
	images = VisitImg(images, htmlNodes)

	scanner := bufio.NewScanner(bytes.NewReader(htmdata))
	scanner.Split(bufio.ScanWords)
	var count int
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		fmt.Errorf("单词分隔浏览发生错误%s", err)
		return 0, len(images)
	}
	resp.Body.Close()
	for _, src := range images {
		fmt.Println(src)
	}
	return count, len(images)

}

//尝试重新连接服务器---学习错误处理
func WaitForServer(url string) error {
	const timeout = 1 * time.Minute
	deadline := time.Now().Add(timeout)
	for trys := 0; time.Now().Before(deadline); trys++ {

		_, err := http.Head(url)
		if err == nil {
			return nil
		}
		log.Printf("server not responding(%s)", err)
		time.Sleep(time.Second << uint(trys)) //指数退避策略
	}
	return fmt.Errorf("server %s failed to respond after %s", url, timeout)
}

//解析html并返回文档中存在的链接---使用匿名函数
func Extract(url string, linktype string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("服务器返回错误,.url:%s,error:%s", url, resp.StatusCode)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("html文档解析发生错误%v", err)
	}
	var links []string
	visitNode := func(node *html.Node) {

		if node.Type == html.ElementNode && node.Data == linktype {
			for _, a := range node.Attr {
				if linktype == "img" {
					if a.Key != "src" {
						continue
					}
				}
				if linktype == "a" {
					if a.Key != "href" {
						continue
					}
				}

				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue
				}
				links = append(links, link.String())
			}
		}
	}
	ForEachNode(doc, visitNode, nil)
	return links, nil
}

//使用goquery解析html文档，下载蝌蚪窝视频
func ExtractOfKeDouWO(url string) []string {
	p, err := goquery.ParseUrl(url)
	if err != nil {
		panic(err)
	}
	var links []string
	item := p.Find("a")
	item.Each(func(index int, element *goquery.Node) {
		var url, videokey string
		for _, v := range element.Attr {
			if v.Key == "href" {
				url = v.Val
			}
			if v.Key == "data-attach-session" {
				videokey = v.Val
			}
		}
		if (videokey != "") && (url != "") {
			videotitle := p.Find("title").Text()
			writeFile("蝌蚪窝视频.txt", videotitle+"    "+url+"\r\n")
			return
		}
		if url != "" {
			temurl, err := url2.ParseRequestURI(url)
			if err != nil {
				return
			}
			if temurl.Host == "www.cao0002.com" {
				links = append(links, url)
			}

		}

	})
	return links
}

//广度遍历，不重复访问节点,爬出相关资源--学习匿名函数
func BreadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}

	}
}

//解析所有连接
func Crawl(url string) []string {
	fmt.Println(url)
	list, err := Extract(url, "a")
	if err != nil {
		log.Print(err)
	}
	return list
}

//闭包问题测试
func ClosureTest() {
	m := make(map[int]func() int)
	for i := 0; i < 4; i++ {
		j := i //这一句很关键，如果直接使用i是不合适的，因为i的值随着循环会迭代更新
		m[i] = func() int {
			fmt.Println(j)
			return j
		}
	}
	for _, v := range m {
		v()
	}

}

//变长函数--多参数函数
func MaxInt(nums ...int) int {
	num := 0
	if nums == nil {
		return 0
	}
	for _, v := range nums {
		if num < v {
			num = v
		}
	}
	return num

}

//下载url内容---简单defer的使用
func Fetch(url string) (name string, n int64, err error) {
	reps, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer reps.Body.Close()
	local := path.Base(reps.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}
	if !strings.Contains(local, ".") {
		local += ".html"
	}
	f, err := os.Create(local)
	if err != nil {
		return "", 0, nil
	}
	n, err = io.Copy(f, reps.Body)
	if closeErr := f.Close(); err != nil {
		err = closeErr
	}
	return local, n, err
}
func writeFile(fileName string, content string) {
	fout, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer fout.Close()
	if err != nil {
		fmt.Println(fileName, err)
		return
	}
	fout.WriteString(content)
}
