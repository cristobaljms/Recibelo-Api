package controllers

import (
	"fmt"
	"net/http"
)

// Index return a blank page
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome convive")
}