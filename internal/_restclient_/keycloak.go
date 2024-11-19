package restclient

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/dumacp/go-driverconsole/internal/utils"
	"github.com/dumacp/go-logs/pkg/logs"
	"golang.org/x/oauth2"
)

func Token(
	user, pass, url, realm, clientID, clientSecret string) (oauth2.TokenSource, *http.Client, error) {

	c := &http.Client{
		Transport: utils.LoadLocalCert(),
	}
	c.Timeout = 30 * time.Second
	ctx_ := context.TODO()
	ctx := context.WithValue(ctx_, oauth2.HTTPClient, c)
	tks, err := TokenSource(ctx,
		user, pass, url, realm, clientID, clientSecret)
	if err != nil {
		logs.LogError.Println(err)
		fmt.Println(err)
		return nil, nil, err
	}
	client := oauth2.NewClient(ctx, tks)
	return tks, client, nil
}

func TokenSource(ctx context.Context, username, password, url, realm, clientid, clientsecret string) (oauth2.TokenSource, error) {

	issuer := fmt.Sprintf("%s/realms/%s", url, realm)
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, err
	}

	config := &oauth2.Config{
		ClientID:     clientid,
		ClientSecret: clientsecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  url,
		Scopes:       []string{oidc.ScopeOpenID},
	}

	tk, err := config.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		return nil, err
	}
	ts := config.TokenSource(ctx, tk)
	return ts, nil

}
