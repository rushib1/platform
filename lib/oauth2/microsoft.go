/*
 * Copyright (c) 2019. Pandranki Global Private Limited
 */

package oauth2

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
	"io/ioutil"
	"net/http"
	"os"
)

// Scopes: OAuth 2.0 scopes provide a way to limit the amount of access that is granted to an access token.
var microsoftOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8000/auth/microsoft/callback",
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	Scopes:       []string{"User.Read", "Contacts.Read"},
	Endpoint:     microsoft.AzureADEndpoint(""),
}

const microsoftEndpoint = "https://graph.microsoft.com"
const microsoftUserInfoEndpoint = microsoftEndpoint + "v1.0/me/"
const oauthMicrosoftUrlAPI = "https://graph.microsoft.com/v1.0/me/?access_token="

func oauthMicrosoftLogin(w http.ResponseWriter, r *http.Request) {

	// Create oauthState cookie
	oauthState := generateStateOauthCookie(w)

	/*
	   AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
	   validate that it matches the the state query parameter on your redirect callback.
	*/
	u := microsoftOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func oauthMicrosoftCallback(w http.ResponseWriter, r *http.Request) {
	// Read oauthState from Cookie
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth microsoft state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := getUserDataFromMicrosoft(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// GetOrCreate User in your db.
	// Redirect or response with a token.
	// More code .....
	fmt.Fprintf(w, "UserInfo: %s\n", data)
}

func getUserDataFromMicrosoft(code string) ([]byte, error) {
	// Use code to get token and get user info from Microsoft.

	token, err := microsoftOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthMicrosoftUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}

func GetUserDataFromMicrosoftUsingAccessToken(accessToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", microsoftUserInfoEndpoint, nil)
	client := &http.Client{}

	if err != nil {
		log.Printf("Get: %s\n", err)
		return nil, err
	}
	// Setting get parameters
	params := req.URL.Query()
	params.Add("access_token", accessToken)

	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Errorln("Unable to process SocialAuth %s", err)
		return nil, err
	} else if resp.StatusCode != 200 {
		log.Errorln("Unable to process SocialAuth %s", http.StatusText(resp.StatusCode))
		return nil, errors.New("Unable to process SocialAuth %s" + http.StatusText(resp.StatusCode))
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
