package main

import (
	"os"
)

func getENV() {
	if len(url) <= 0 {
		if len(os.Getenv("URL_DEVICES")) > 0 {
			url = os.Getenv("URL_DEVICES")
		} else {
			url = url_
		}
	}
}
