package main

import (
	"github.com/sirupsen/logrus"
	"html"
	"net/http"
)

func main() {
	logrus.Printf("Starting k8s-ces-setup...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logrus.Printf("Hello, %q", html.EscapeString(r.URL.Path))
	})

	logrus.Fatal(http.ListenAndServe(":8888", nil))
}
