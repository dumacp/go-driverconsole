package main

import (
	"os"
)

func getENV() {
	if len(url) <= 0 {
		if len(os.Getenv("APPFARE_URL")) > 0 {
			url = os.Getenv("APPFARE_URL")
		} else {
			url = url_
		}
	}
}
