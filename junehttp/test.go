package junehttp

import (
	"fmt"
	"net/http"
)

func Servertest() {
	http.HandleFunc("/postpage", func(w http.ResponseWriter, r *http.Request) {
		//接受post请求，然后打印表单中key和value字段的值
		if r.Method == "POST" {
			var (
				key   string = r.PostFormValue("key")
				value string = r.PostFormValue("value")
			)
			fmt.Printf("key is  : %s\n", key)
			fmt.Printf("value is: %s\n", value)
		}
		for _, que := range r.URL.Query() {
			fmt.Println(que)
		}
		values := r.PostForm["key"]
		for _, v := range values {
			fmt.Println(v)
		}

	})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func Clienttest() {
	DoHttpRequest(
		"http://localhost:8000/postpage",
		"POST",
		"key=fuchongjun&value=18",
		"",
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		nil, map[string]string{"para": "444"}, nil,
		30)
}
