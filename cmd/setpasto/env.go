package main

import (
	"os"
)

func getENV() {
	if len(url) <= 0 {
		if len(os.Getenv("KEYCLOAK_URL_DEVICES")) > 0 {
			url = os.Getenv("KEYCLOAK_URL_DEVICES")
		} else {
			url = url_
		}
	}
}
