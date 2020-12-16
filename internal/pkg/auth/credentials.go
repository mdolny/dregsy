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

package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

//
type Refresher interface {
	Refresh(creds *Credentials) error
}

//
func NewCredentialsFromBasic(username, password string, json bool) (*Credentials, error) {
	return &Credentials{
		username: username,
		password: password,
		jsonAuth: json}, nil
}

//
func NewCredentialsFromToken(token string) (*Credentials, error) {
	return &Credentials{token: NewToken(token)}, nil
}

//
func NewCredentialsFromAuth(auth string) (*Credentials, error) {

	data, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return nil, err
	}

	ret := &Credentials{}

	crd := &jsonCreds{}
	if err := json.Unmarshal(data, crd); err != nil {
		ret.jsonAuth = false
		parts := strings.SplitN(string(data), ":", 2)
		ret.username = parts[0]
		if len(parts) > 1 {
			ret.password = parts[1]
		}
	} else {
		ret.jsonAuth = true
		ret.username = crd.User
		ret.password = crd.Pass
	}

	return ret, nil
}

//
type jsonCreds struct {
	User string `json:"username"`
	Pass string `json:"password"`
}

//
type Credentials struct {
	//
	username string
	password string
	jsonAuth bool
	//
	token     *Token
	refresher Refresher
}

//
func (c *Credentials) Username() string {
	return c.username
}

//
func (c *Credentials) Password() string {
	return c.password
}

//
func (c *Credentials) Auth() string {

	if c.username == "" && c.password == "" {
		return ""
	}

	template := "%s:%s"
	if c.jsonAuth {
		template = `{"username": "%s", "password": "%s"}`
	}

	return base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf(template, c.username, c.password)))
}

//
func (c *Credentials) Token() *Token {
	return c.token
}

//
func (c *Credentials) SetToken(t *Token) {
	c.token = t
}

//
func (c *Credentials) SetRefresher(r Refresher) {
	c.refresher = r
}

//
func (c *Credentials) Refresh() error {
	if c.refresher == nil {
		return nil
	}
	return c.refresher.Refresh(c)
}
