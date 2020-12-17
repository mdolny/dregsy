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
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/xelalexv/dregsy/internal/pkg/auth"
)

//
type ListSourceType string

const (
	Catalog   ListSourceType = "catalog"
	DockerHub                = "dockerhub"
	Index                    = "index"
)

//
func IsValidListSourceType(t ListSourceType) bool {
	switch t {
	case Catalog, DockerHub, Index:
		return true
	}
	return false
}

//
type ListSource interface {
	Ping() error
	Retrieve() ([]string, error)
}

//
func NewRepoList(registry string, typ ListSourceType, creds *auth.Credentials) (
	*RepoList, error) {

	list := &RepoList{registry: registry}

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

	// DockerHub does not expose the registry catalog API, but separate APIs for
	// listing and searching. These APIs use tokens that are different from the
	// one used for normal registry actions, so we clone the credentials for list
	// use. For listing via catalog API, we can use the same credentials as for
	// push & pull.
	listCreds := creds
	if server == "registry.hub.docker.com" {
		var err error
		listCreds, err = auth.NewCredentialsFromBasic(
			creds.Username(), creds.Password())
		if err != nil {
			return nil, err
		}
		if typ != DockerHub && typ != Index {
			return nil, fmt.Errorf(
				"DockerHub only supports list types '%s' and '%s'",
				DockerHub, Index)
		}
	}

	switch typ {

	case DockerHub:
		list.source = newDockerhub(reg, insecure, listCreds)

	case Index:
		list.source = newIndex(reg, listCreds.Username(), insecure, listCreds)

	case Catalog, "":
		list.source = newCatalog(
			reg, insecure, strings.HasSuffix(server, ".gcr.io"), listCreds)

	default:
		return nil, fmt.Errorf("invalid list source type '%s'", typ)
	}

	return list, nil
}

//
type RepoList struct {
	registry string
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

	var err error
	if l.repos, err = l.source.Retrieve(); err != nil {
		return nil, err
	}

	l.expiry = time.Now().Add(10 * time.Minute) // FIXME: parameterize

	return l.repos, nil
}
