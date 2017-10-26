package junehttp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

func DoGet(url string, contenttype string) ([]byte, error) {
	client := http.DefaultClient
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if contenttype != "" {
		request.Header.Add("Content-Type", contenttype)

	}
	respons, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer respons.Body.Close()
	return ioutil.ReadAll(respons.Body)
}
func DoDelete(url string, contenttype string) ([]byte, error) {
	client := http.DefaultClient
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	if contenttype != "" {
		request.Header.Add("Content-Type", contenttype)

	}
	respons, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer respons.Body.Close()
	return ioutil.ReadAll(respons.Body)
}
func DoFormPost(url string, values url.Values) ([]byte, error) {

	//把post表单发送给目标服务器
	res, err := http.PostForm(url, values)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)

}

//method:POSt,GET,DELETE,PUT
//data:如果是form提交，格式为"key1=value1&key2=value2"，也可以设置formdata参数,原理一样
//headers:from表单post Content-Type：application/x-www-form-urlencoded,
//json提交设置为application/json
//proxy 代理服务器，可选
//timeoutSeconds 不设置可以设为0
//formdata:格式化数据格式为"key1=value1&key2=value2",一定要设置Content-Type：application/x-www-form-urlencoded
//params:在请求的URL后面设置请求参数
//注意：如果以form形式提交，一定要设置头Content-Type：application/x-www-form-urlencoded
func DoHttpRequest(uri, method, bodydata, proxy string, headers, formdata, params, cookies map[string]string, timeoutSeconds int) ([]byte, error) {

	//设置请求数据
	data := ""
	if bodydata != "" {
		data = bodydata
	}
	if formdata != nil {
		var req http.Request
		req.ParseForm()
		for key, val := range formdata {
			req.Form.Add(key, val)
		}
		data = strings.TrimSpace(req.Form.Encode())
	}

	request, err := http.NewRequest(method, uri, strings.NewReader(data))
	//设置url后面的请求参数，其实就是get请求参数
	if params != nil {
		q := request.URL.Query()
		for key, val := range params {
			q.Add(key, val)
		}
		request.URL.RawQuery = q.Encode()
	}
	//设置header
	if headers != nil {
		for key, val := range headers {
			request.Header.Set(key, val)
		}
	}
	//设置cookie
	if cookies != nil {
		for key, val := range cookies {
			cookie := &http.Cookie{Name: key, Value: val, HttpOnly: false}
			request.AddCookie(cookie)
		}
	}
	//设置代理服务器
	transport := &http.Transport{} //这一步很重要，设置了很多默认传输参数
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
	}

	client := &http.Client{ //多次请求，并发请求不要用默认的defaultclient,单次请求可以使用下
		Transport: transport,
	}
	//设置请求超时时间
	if timeoutSeconds > 0 {
		client.Timeout = time.Duration(timeoutSeconds) * time.Second
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Printf("request error:%s", err)
		return nil, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("request error:%s", err)
		return nil, err
	}

	return result, err
}

//body提交二进制数据
func DoBytesPost(url string, data []byte, headers map[string]string) ([]byte, error) {

	body := bytes.NewReader(data)
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Println("http.NewRequest,[err=%s][url=%s]", err, url)
		return []byte(""), err
	}
	if headers != nil {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}
	var resp *http.Response
	transport := &http.Transport{}
	httpclient := &http.Client{
		Transport: transport,
	}
	resp, err = httpclient.Do(request)
	if err != nil {
		log.Println("http.Do failed,[err=%s][url=%s]", err, url)
		return []byte(""), err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("http.Do failed,[err=%s][url=%s]", err, url)
	}
	return b, err
}

//构建文件上传请求体：request
//uri:上传路径
//params:表单键值参数,除了文件其他表单参数
//paramName:上传控件的文件域名字，即服务端input标签的name
//filepath文件路径
func NewfileUploadRequest(uri string, params map[string]string, paramName, filepath string) (*http.Request, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//构建请求体，新建body指针
	body := &bytes.Buffer{}
	//新建body数据writer
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, path.Base(filepath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", uri, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request, err
}

//模拟文件上传和NewfileUploadRequest功能差不多
func PostFile(fieldname string, filename string, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile(fieldname, path.Base(filename))
	fmt.Println(path.Base(filename))
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return nil
}

func NewfileUploadRequest_test() {
	//路径要用这种格式，用反斜杠
	filepath := `D:/IMG_1618.JPG`
	m := map[string]string{"": ""}
	req, err := NewfileUploadRequest("http://localhost:60230/Default2", m, "FileUpload1", filepath)
	if err != nil {
		log.Println(err.Error())
	}
	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(b))
}

func PostFile_test() {
	target_url := "http://localhost:12345/hello"
	filename := "D:/IMG_1618.JPG"
	err := PostFile("userfile", filename, target_url)
	fmt.Println(err)
}
