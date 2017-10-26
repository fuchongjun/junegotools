package study

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type dollars float32

func (d dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

type database map[string]dollars

func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/list":
		for item, price := range db {
			fmt.Fprintf(w, "%s:%s\n", item, price)
		}
	case "/price":
		item := req.URL.Query().Get("item")
		price, ok := db[item]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "no such item %q\n", item)
			return
		}
		fmt.Fprintf(w, "%s\n", price)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such page:%s\n", req.URL)
	}

}
func HandlerTest() {
	db := database{"shose": 50, "socks": 5}
	log.Fatal(http.ListenAndServe("localhost:8000", db))
}

//=======拆分处理=======
type dbkey []string

func (db database) list(w http.ResponseWriter, req *http.Request) {
	var dbkeys dbkey
	for item, _ := range db {
		dbkeys = append(dbkeys, item)
	}
	sort.Sort(dbkeys)
	for _, v := range dbkeys {
		fmt.Fprintf(w, "%s:%s\n", v, db[v])
	}
}
func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item %q\n", item)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}
func (db database) add(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price, err := strconv.ParseFloat(r.URL.Query().Get("price"), 32)
	if err != nil {
		http.Error(w, "添加产品价格类型错误", http.StatusNotAcceptable)
		return
	}
	if _, ok := db[item]; ok {
		http.Error(w, "该类型产品已经添加", http.StatusNotAcceptable)
		return
	}
	db[item] = dollars(price)
	fmt.Fprintf(w, "添加商品成功:%s:%s", item, db[item])
}

func (db dbkey) Len() int {
	return len(db)
}
func (db dbkey) Less(i, j int) bool {
	return db[i] < db[j]
}
func (db dbkey) Swap(i, j int) {
	db[i], db[j] = db[j], db[i]
}

func HandlerTest2() {
	db := database{"shose": 50, "socks": 5}
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/add", db.add)
	err:=http.ListenAndServe("localhost:8080", nil)
	if err!=nil {
		log.Fatal(err)
	}

}
