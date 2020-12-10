/*
	Copyright 2020 Alexander Vollschwitz <xelalex@gmx.net>

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	  http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

//
type DHRepoList struct {
	Items []DHRepoDescriptor `json:"results",required`
}

//
type DHRepoDescriptor struct {
	User        string `json:"user",required`
	Name        string `json:"name",required`
	Namespace   string `json:"namespace",required`
	RepoType    string `json:"repository_type",required`
	Description string `json:"description",omitempty`
	IsPrivate   bool   `json:"is_private",omitempty`
}

//
func newDockerhub(reg string, insecure bool, auth *Auth) ListSource {
	return &dockerhub{auth: auth}
}

//
type dockerhub struct {
	auth *Auth
}

//
func (d *dockerhub) Retrieve() ([]string, error) {

	token, err := d.getToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		"https://hub.docker.com/v2/repositories/%s/?page_size=1000",
		d.auth.Username)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("JWT %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var list DHRepoList
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}

	ret := make([]string, 0, len(list.Items))
	for _, r := range list.Items {
		ret = append(ret, fmt.Sprintf("%s/%s", r.User, r.Name))
	}

	return ret, nil
}

//
func (d *dockerhub) Ping() error {
	_, err := d.getToken()
	return err
}

//
func (d *dockerhub) getToken() (string, error) {

	auth := url.Values{
		"username": {d.auth.Username},
		"password": {d.auth.Password},
	}

	resp, err := http.PostForm("https://hub.docker.com/v2/users/login/", auth)
	if err != nil {
		return "", err
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	token, ok := res["token"].(string)
	if !ok {
		return "", fmt.Errorf("received token is not a string")
	}
	return token, nil
}
