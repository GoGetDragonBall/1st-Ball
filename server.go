package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func test(w http.ResponseWriter, r *http.Request) {
	//Parsing HTML
	t, err := template.ParseFiles("html/test.html")
	if err != nil {
		fmt.Println(err)
	}

	items := struct {
		Name string
		City string
	}{
		Name: "MyName",
		City: "MyCity",
	}

	t.Execute(w, items)
}

type Data struct {
	Name string
}

func testInput(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	s := r.FormValue("test")
	fmt.Println(s)
	return
}
func testInput2(w http.ResponseWriter, r *http.Request) {
	fmt.Println("input2" + r.Method)
	fmt.Println(r.Method)
	r.ParseForm()

	data := Data{"ajax 리턴 데이터! 입력된 값은 " + r.FormValue("test")}
	w.Header().Set("Content-type", "application/json")
	err := json.NewEncoder(w).Encode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var mux map[string]func(http.ResponseWriter, *http.Request)

func defineMux() {
	mux["/"] = hello
	mux["/test"] = test
	mux["/test/input"] = testInput
	mux["/test/input2"] = testInput2
}

func main() {
	server := &http.Server{
		Addr:           ":8080",
		Handler:        &myHandler{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	defineMux()

	server.ListenAndServe()
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.String())
	str := strings.Split(r.URL.String(), "?")[0]
	if h, ok := mux[str]; ok {

		h(w, r)
		return
	}
	io.WriteString(w, "My server: "+r.URL.String())

}
