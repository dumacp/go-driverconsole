package utils

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dumacp/go-logs/pkg/logs"
)

const (
	localCertDir = "/usr/local/share/ca-certificates/"
)

var ErrorBadRequest = errors.New("bad request")

func Post(client *http.Client,
	url, username, password string,
	jsonStr []byte) ([]byte, int, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	if len(username) > 0 || len(password) > 0 {
		req.SetBasicAuth(username, password)
	}

	if client == nil {
		tr := LoadLocalCert()
		client = &http.Client{Transport: tr}
		client.Timeout = 30 * time.Second
	}

	logs.LogBuild.Printf("Post request: %+v", req)

	var resp *http.Response
	rangex := make([]int, 3)
	for range rangex {
		resp, err = client.Do(req)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}
	if err != nil {
		return nil, 0, fmt.Errorf("post error, err: request: %+v, err: %w", req, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		if err != nil {
			return nil, resp.StatusCode, err
		}
		if resp.StatusCode == 400 {
			return nil, resp.StatusCode, fmt.Errorf("%w, StatusCode: %d, resp: %s, req: %+v",
				ErrorBadRequest, resp.StatusCode, body, req)
		}
		return nil, resp.StatusCode, fmt.Errorf("StatusCode: %d, resp: %s, req: %+v", resp.StatusCode, body, req)
	}
	return body, resp.StatusCode, nil
}

func Get(client *http.Client,
	url, username, password string, jsonStr []byte) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	if len(username) > 0 && len(password) > 0 {
		req.SetBasicAuth(username, password)
	}

	if client == nil {
		tr := LoadLocalCert()
		client = &http.Client{Transport: tr}
		client.Timeout = 30 * time.Second
	}

	logs.LogBuild.Printf("Get request: %+v", req)
	var resp *http.Response
	rangex := make([]int, 3)
	for range rangex {
		resp, err = client.Do(req)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		if err != nil {
			return nil, resp.StatusCode, err
		}
		return nil, resp.StatusCode, fmt.Errorf("StatusCode: %d, resp: %s, req: %s", resp.StatusCode, body, req.URL)
	}
	return body, resp.StatusCode, nil
}

func LoadLocalCert() *http.Transport {

	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Read in the cert file
	certs, err := os.ReadDir(localCertDir)
	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", localCertDir, err)
	}
	for _, cert := range certs {
		file, err := os.ReadFile(localCertDir + cert.Name())
		if err != nil {
			log.Fatalf("Failed to append %q to RootCAs: %v", cert, err)
		}
		// Append our cert to the system pool
		if ok := rootCAs.AppendCertsFromPEM(file); !ok {
			log.Println("No certs appended, using system certs only")
		}
	}

	// Trust the augmented cert pool in our client
	config := &tls.Config{
		//InsecureSkipVerify: *insecure,
		RootCAs: rootCAs,
	}
	tr := &http.Transport{
		TLSClientConfig: config,
		Dial: (&net.Dialer{
			Timeout: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		// Dial: (&net.Dialer{
		// 	Timeout:   30 * time.Second,
		// 	KeepAlive: 60 * time.Second,
		// }).Dial,
		// TLSHandshakeTimeout:   10 * time.Second,
		// ResponseHeaderTimeout: 10 * time.Second,
		// ExpectContinueTimeout: 3 * time.Second,
	}
	return tr
}
