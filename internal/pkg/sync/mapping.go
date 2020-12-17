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

package sync

import (
	"fmt"
	"regexp"
	"strings"
)

//
const RegexpPrefix = "regex:"

//
type Mapping struct {
	From string   `yaml:"from"`
	To   string   `yaml:"to"`
	Tags []string `yaml:"tags"`
	//
	fromFilter *regexp.Regexp
}

//
func (m *Mapping) validate() error {

	if m == nil {
		return fmt.Errorf("mapping is nil")
	}

	if m.From == "" {
		return fmt.Errorf("mapping without 'From' path")
	}

	if m.isRegexp() {
		regex := m.From[len(RegexpPrefix):]
		var err error
		if m.fromFilter, err = compileRegex(regex); err != nil {
			return fmt.Errorf("invalid regular expression '%s': %v", regex, err)
		}

	} else {
		if m.To == "" {
			m.To = m.From
		}
		m.From = normalizePath(m.From)
	}

	m.To = normalizePath(m.To)

	return nil
}

//
func (m *Mapping) filterRepos(repos []string) []string {

	if m.isRegexp() {
		ret := make([]string, 0, len(repos))
		for _, r := range repos {
			if m.fromFilter.MatchString(r) {
				ret = append(ret, normalizePath(r))
			}
		}
		return ret
	}

	return repos
}

//
func (m *Mapping) isRegexp() bool {
	return strings.HasPrefix(m.From, RegexpPrefix)
}

//
func compileRegex(v string) (*regexp.Regexp, error) {
	if !strings.HasPrefix(v, "^") {
		v = fmt.Sprintf("^%s", v)
	}
	if !strings.HasSuffix(v, "$") {
		v = fmt.Sprintf("%s$", v)
	}
	return regexp.Compile(v)
}

//
func normalizePath(p string) string {
	if strings.HasPrefix(p, "/") {
		return p
	}
	return "/" + p
}
