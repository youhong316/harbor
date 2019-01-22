// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"github.com/goharbor/harbor/src/common/secret"
	"net/http"
)

// Authenticator defines Authenticate function to authenticate requests
type Authenticator interface {
	// Authenticate the request, if there is no error, the bool value
	// determines whether the request is authenticated or not
	Authenticate(req *http.Request) (bool, error)
}

type secretAuthenticator struct {
	secrets map[string]string
}

// NewSecretAuthenticator returns an instance of secretAuthenticator
func NewSecretAuthenticator(secrets map[string]string) Authenticator {
	return &secretAuthenticator{
		secrets: secrets,
	}
}

// Authenticate the request according the secret
func (s *secretAuthenticator) Authenticate(req *http.Request) (bool, error) {
	if len(s.secrets) == 0 {
		return true, nil
	}
	reqSecret := secret.FromRequest(req)

	for _, v := range s.secrets {
		if reqSecret == v {
			return true, nil
		}
	}

	return false, nil
}
