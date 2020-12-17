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
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/xelalexv/dregsy/internal/pkg/auth"
)

//
type ListSource interface {
	Ping() error
	Retrieve() ([]string, error)
}

//
func NewRepoList(registry, filter string, creds *auth.Credentials) (
	*RepoList, error) {

	rxf, err := CompileRegex(filter)
	if err != nil {
		return nil, err
	}

	list := &RepoList{registry: registry, filter: rxf}

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
		// DockerHub does not expose the registry catalog API, but separate
		// APIs for listing and searching. These APIs use tokens that are
		// different from the one used for normal registry actions, so we
		// clone the credentials for list use.
		//
		//list.source = newIndex(reg, creds.Username(), insecure, creds)
		if credsCopy, err := auth.NewCredentialsFromBasic(
			creds.Username(), creds.Password()); err != nil {
			return nil, err
		} else {
			list.source = newDockerhub(reg, insecure, credsCopy)
		}

	default:
		// for listing via catalog API, we can use the same credentials as for
		// push & pull
		list.source = newCatalog(reg, insecure, creds)
	}

	return list, nil
}

//
type RepoList struct {
	registry string
	filter   *regexp.Regexp
	repos    []string
	expiry   time.Time
	source   ListSource
}

//
func (l *RepoList) Get() ([]string, error) {

	if time.Now().Before(l.expiry) {
		log.Debug("list still valid, reusing")
		return l.repos, nil
	}

	log.Debug("retrieving list")

	raw, err := l.source.Retrieve()
	if err != nil {
		return nil, err
	}

	l.expiry = time.Now().Add(10 * time.Minute) // FIXME: parameterize
	l.repos = make([]string, 0, len(raw))
	for _, r := range raw {
		if l.filter.MatchString(r) {
			l.repos = append(l.repos, r)
		}
	}

	return l.repos, nil
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
