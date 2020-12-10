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
	"fmt"
	"regexp"
	"strings"
)

//
type Auth struct {
	Username string
	Password string
	Token    string
}

//
func NewRepoList(registry, filter string, auth *Auth) (*RepoList, error) {

	rxf, err := CompileRegex(filter)
	if err != nil {
		return nil, err
	}

	list := &RepoList{registry: registry, filter: rxf, auth: auth}

	insecure := false
	reg := ""

	if strings.HasPrefix(registry, "http://") {
		insecure = true
		reg = registry[7:]
	} else if strings.HasPrefix(registry, "https://") {
		reg = registry[8:]
	} else {
		reg = registry
	}

	server := strings.SplitN(reg, ":", 2)[0]

	switch server {

	case "registry.hub.docker.com":
		//list.source = newIndex(reg, auth.Username, insecure, auth)
		list.source = newDockerhub(reg, insecure, auth)

	default:
		list.source = newCatalog(reg, insecure, auth)
	}

	return list, nil
}

//
type RepoList struct {
	registry string
	auth     *Auth
	filter   *regexp.Regexp
	repos    []string
	source   ListSource
}

//
func (l *RepoList) Get() ([]string, error) {

	raw, err := l.source.Retrieve()
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0, len(raw))
	for _, r := range raw {
		if l.filter.MatchString(r) {
			ret = append(ret, r)
		}
	}

	return ret, nil
}

//
type ListSource interface {
	Ping() error
	Retrieve() ([]string, error)
}

//
func CompileRegex(v string) (*regexp.Regexp, error) {
	if !strings.HasPrefix(v, "^") {
		v = fmt.Sprintf("^%s", v)
	}
	if !strings.HasSuffix(v, "$") {
		v = fmt.Sprintf("%s$", v)
	}
	return regexp.Compile(v)
}
