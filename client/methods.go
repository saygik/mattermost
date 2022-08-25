/*
  Copyright 2021-2022 Davide Madrisan <davide.madrisan@gmail.com>

  Licensed under the Mozilla Public License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      https://www.mozilla.org/en-US/MPL/2.0/

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

// Package mattermost implements the API v4 calls to Mattemost.
package mattermost

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/saygik/go-glpi-to-matt-test/config"
	"io"
	"io/ioutil"
	"net/http"
)

// queryAPIv4 makes a query to Mattermost using its REST API v4.
func queryAPIv4(method, endpoint string, payload io.Reader, opts config.Options) (interface{}, error) {
	baseURL, err := getURL()
	if err != nil {
		return nil, err
	}

	accessToken, err := getAccessToken()
	if err != nil {
		return nil, err
	}

	var bearer = forgeBearerAuthentication(accessToken)
	var url = forgeAPIv4URL(baseURL, endpoint)

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Accept", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: opts.SkipTLSVerify,
		},
	}

	client := &http.Client{
		Timeout:   opts.ConnectionTimeout,
		Transport: tr,
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		statusCodeText := http.StatusText(response.StatusCode)
		return nil, fmt.Errorf("the HTTP query to %s has ended with a %d (\"%s\") code",
			url, response.StatusCode, statusCodeText)
	}

	// Read body
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var data interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return nil, err
	}

	return data, nil
}

// Get makes a query of type GET to Mattermost.
func Get(endpoint string) (interface{}, error) {
	var opts = config.Options{}
	response, err := queryAPIv4(http.MethodGet, endpoint, nil, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Post makes a query of type POST to Mattermost.
func Post1(endpoint string, payload io.Reader, opts config.Options) (interface{}, error) {
	response, err := queryAPIv4(http.MethodPost, endpoint, payload, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}