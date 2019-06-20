package main

import (
	"net/http"
	"fmt"
	"crypto/md5"
	"io"
	"time"
	"encoding/hex"
)

type Answer struct {
	md5 string
	status string
	url string
	id string
}

type Message struct {
	ans Answer
	typ string
}

var DataRequests chan Message = make(chan Message) 
var DataResponses chan Message = make(chan Message)

func DataStorage() {
	Jobs := make(map[string]Answer)
	for {
		Current := <- DataRequests
		if Current.typ == "post" {
			Jobs[Current.ans.id] = Current.ans
			if Current.ans.status == "running" {
				go Md5Counter(Current.ans.url, Current.ans.id)
				fmt.Println("Job started, url:", Current.ans.url, ", id:",Current.ans.id)
			} 
		} else if Current.typ == "get" {
			fmt.Println("Job with id", Current.ans.id, "requested")
			res, ok := Jobs[Current.ans.id]
			if !ok {
				DataResponses <- Message{ans: Answer{md5: "", status: "not exists", url: "", id: Current.ans.id}, typ: "resp"}
			} else {
				DataResponses <- Message{ans: res, typ: "resp"}
			}
		}
	}
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.Form.Get("url")
	hasher := md5.New()
	io.WriteString(hasher, url)
	io.WriteString(hasher, time.Now().String())
	id := fmt.Sprintf("%x",  hasher.Sum(nil))
	DataRequests <- Message{ans: Answer{md5: "", status: "running", url: url, id: id}, typ: "post"}
	io.WriteString(w, "{\"id\":\"" + id + "\"}\n")
	go Md5Counter(url, id)	
}

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form.Get("id")
	DataRequests <- Message{ans: Answer{md5:"", status:"", url:"", id:id}, typ: "get"}
	Answer := <- DataResponses
	io.WriteString(w, "{\"md5\":\"" + Answer.ans.md5 + "\",\"status\":" + Answer.ans.status + "\"}\n")
}

func Md5Counter(url string, id string) {
	
	resp, err := http.Get(url)
    if err != nil {
		DataRequests <- Message{ans: Answer{md5: "", status:"error downloading a file", url:url, id:id}, typ: "post"}
		return
    }
    defer resp.Body.Close()

	hasher := md5.New()
	
    _, err = io.Copy(hasher, resp.Body)
    if err != nil {
		DataRequests <- Message{ans: Answer{md5: "", status:"error writing a file", url:url, id:id}, typ: "post"}
		return
	}
	DataRequests <- Message{ans: Answer{md5: hex.EncodeToString(hasher.Sum(nil)), status:"finished", url:url, id:id}, typ: "post"}
}


func main() {

	http.HandleFunc("/submit", SubmitHandler)
	http.HandleFunc("/check", CheckHandler)
	go DataStorage()
	http.ListenAndServe(":8000", nil)
}