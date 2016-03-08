package main

import (
  "fmt"
  "net/http"
  "log"

  "github.com/zenazn/goji"
  "github.com/zenazn/goji/web"
)

func main() {
  http.Handle("/", http.FileServer(http.Dir("public")))
  fmt.Println("Server started on port 8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
