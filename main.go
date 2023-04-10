package main

import (
  "fmt"
  "net/http"
  "os"
)

func main(){
  http.HandleFunc("/", handler)
  http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request){
  msg := os.Getenv("MSG")
  fmt.Fprintf(w, msg)
}