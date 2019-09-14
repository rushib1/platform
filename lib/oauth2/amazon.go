/*
 * Copyright (c) 2019. Pandranki Global Private Limited
 */

package oauth2

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const amazonApiEndpoint = "https://api.amazon.com"
const amazonUserInfoEndpoint = amazonApiEndpoint + "/user/profile"

func GetUserDataFromAmazonUsingAccessToken(accessToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", amazonUserInfoEndpoint, nil)
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
