package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/dumacp/go-driverconsole/internal/utils"
	"golang.org/x/oauth2"
)

const (
	kcRealm        = "DEVICES"
	kcClientID     = "devices2"
	kcClientSecret = "b73479a3-225b-4b96-ad65-22edd82623a3"
	kcBaseURL      = "https://fleet.nebulae.com.co/auth"
	apiBaseURL     = "https://fleet.nebulae.com.co"
)

func main() {
	deviceID := "NE-RCXL-0098"
	doc := "1007185602"

	c := &http.Client{Transport: utils.LoadLocalCert(), Timeout: 60 * time.Second}
	ctx := context.WithValue(context.TODO(), oauth2.HTTPClient, c)
	issuer := fmt.Sprintf("%s/realms/%s", kcBaseURL, kcRealm)
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		fmt.Printf("ERROR oidc: %v\n", err)
		return
	}
	config := &oauth2.Config{
		ClientID: kcClientID, ClientSecret: kcClientSecret,
		Endpoint: provider.Endpoint(), Scopes: []string{oidc.ScopeOpenID},
	}
	tk, err := config.PasswordCredentialsToken(ctx, deviceID, deviceID)
	if err != nil {
		fmt.Printf("ERROR token: %v\n", err)
		return
	}
	client := config.Client(ctx, tk)

	url := fmt.Sprintf("%s/api/external-system-gateway/rest/driver-daily-services/%s?page=0&count=10", apiBaseURL, doc)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("ERROR GET: %v\n", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	b, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(b))
}
