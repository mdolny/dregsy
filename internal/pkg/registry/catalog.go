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
	"context"
	"fmt"

	"golang.org/x/oauth2"

	gcrauthn "github.com/google/go-containerregistry/pkg/authn"
	gcrname "github.com/google/go-containerregistry/pkg/name"
	gcrremote "github.com/google/go-containerregistry/pkg/v1/remote"

	"github.com/xelalexv/dregsy/internal/pkg/auth"
)

//
func newCatalog(reg string, insecure bool, creds *auth.Credentials) ListSource {
	return &catalog{
		registry: reg,
		conf: &oauth2.Config{
			ClientID: creds.Username(),
			Endpoint: oauth2.Endpoint{
				TokenURL: fmt.Sprintf("https://%s/token", reg),
			},
		},
		creds: creds,
	}
}

//
type catalog struct {
	registry string
	conf     *oauth2.Config
	creds    *auth.Credentials
}

//
func (c *catalog) Retrieve() ([]string, error) {

	reg, err := gcrname.NewRegistry(c.registry)
	if err != nil {
		return nil, err
	}

	ret, err := gcrremote.CatalogPage(reg, "", 100,
		gcrremote.WithAuth(&gcrauthn.Basic{
			Username: c.creds.Username(),
			Password: c.creds.Password(),
		}))
	if err != nil {
		return nil, err
	}

	return ret, nil
}

//
func (c *catalog) Ping() error {
	_, err := c.conf.PasswordCredentialsToken(
		context.TODO(), c.creds.Username(), c.creds.Password())
	return err
}
