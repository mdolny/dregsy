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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/registry"
)

//
func newIndex(reg, filter string, insecure bool, auth *Auth) ListSource {

	ret := &index{filter: filter}

	ret.auth = &types.AuthConfig{
		Username:      auth.Username,
		Password:      auth.Password,
		RegistryToken: auth.Token,
	}
	ret.opts = &registry.ServiceOptions{}

	if insecure {
		ret.auth.ServerAddress = fmt.Sprintf("http://%s", reg)
	} else {
		ret.auth.ServerAddress = fmt.Sprintf("https://%s", reg)
	}

	return ret
}

//
type index struct {
	opts   *registry.ServiceOptions
	auth   *types.AuthConfig
	filter string
}

//
func (i *index) Retrieve() ([]string, error) {

	svc, err := registry.NewService(*i.opts)
	if err != nil {
		return nil, err
	}

	res, err := svc.Search(context.TODO(), i.filter, 100, i.auth, "dregsy", nil)
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0, res.NumResults)
	for _, r := range res.Results {
		ret = append(ret, r.Name)
	}

	return ret, nil
}

//
func (i *index) Ping() error {
	svc, err := registry.NewService(*i.opts)
	if err != nil {
		return err
	}
	if _, _, err := svc.Auth(context.TODO(), i.auth, "dregsy"); err != nil {
		return err
	}
	return nil
}
