/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trustbloc/edge-core/pkg/storage/memstore"

	"github.com/trustbloc/hub-auth/pkg/restapi/operation"
)

func TestController_New(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config, cleanup := config()
		defer cleanup()

		controller, err := New(config)
		require.NoError(t, err)
		require.NotNil(t, controller)
	})

	t.Run("error if operations cannot start", func(t *testing.T) {
		config, cleanup := config()
		defer cleanup()
		config.OIDCProviderURL = "BadURL"

		_, err := New(config)
		require.Error(t, err)
	})
}

func TestController_GetOperations(t *testing.T) {
	config, cleanup := config()
	defer cleanup()

	controller, err := New(config)
	require.NoError(t, err)
	require.NotNil(t, controller)

	ops := controller.GetOperations()
	require.Equal(t, 2, len(ops))
}

func config() (*operation.Config, func()) {
	path, cleanup := newTestOIDCProvider()

	return &operation.Config{
		OIDCProviderURL: path,
		Provider:        memstore.NewProvider(),
	}, cleanup
}

func newTestOIDCProvider() (string, func()) {
	h := &testOIDCProvider{}
	srv := httptest.NewServer(h)
	h.baseURL = srv.URL

	return srv.URL, srv.Close
}

type oidcConfigJSON struct {
	Issuer      string   `json:"issuer"`
	AuthURL     string   `json:"authorization_endpoint"`
	TokenURL    string   `json:"token_endpoint"`
	JWKSURL     string   `json:"jwks_uri"`
	UserInfoURL string   `json:"userinfo_endpoint"`
	Algorithms  []string `json:"id_token_signing_alg_values_supported"`
}

type testOIDCProvider struct {
	baseURL string
}

func (t *testOIDCProvider) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	response, err := json.Marshal(&oidcConfigJSON{
		Issuer:      t.baseURL,
		AuthURL:     fmt.Sprintf("%s/oauth2/auth", t.baseURL),
		TokenURL:    fmt.Sprintf("%s/oauth2/token", t.baseURL),
		JWKSURL:     fmt.Sprintf("%s/oauth2/certs", t.baseURL),
		UserInfoURL: fmt.Sprintf("%s/oauth2/userinfo", t.baseURL),
		Algorithms:  []string{"RS256"},
	})
	if err != nil {
		panic(err)
	}

	_, err = w.Write(response)
	if err != nil {
		panic(err)
	}
}