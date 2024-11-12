package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/asecurityteam/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPDefaultComponent(t *testing.T) {
	cmp := &DefaultComponent{}
	conf := cmp.Settings()
	tr, err := cmp.New(context.Background(), conf)
	require.Nil(t, err)
	require.NotNil(t, tr)
}

func TestHTTP(t *testing.T) {
	src := settings.NewMapSource(map[string]interface{}{
		"httpclient": map[string]interface{}{
			"type": "DEFAULT",
		},
	})
	tr, err := New(context.Background(), src)
	require.Nil(t, err)
	require.NotNil(t, tr)

	src = settings.NewMapSource(map[string]interface{}{
		"httpclient": map[string]interface{}{
			"type": "SMART",
		},
	})
	_, err = New(context.Background(), src)
	require.NotNil(t, err) //must bomb out on attempt to create smart client

	src = settings.NewMapSource(map[string]interface{}{
		"httpclient": map[string]interface{}{
			"type": "MISSING",
		},
	})
	_, err = New(context.Background(), src)
	require.NotNil(t, err)
}

/*
Given a default http client with a default Content-Type of application/json
When a request is sent without a Content-Type header defined
Then the request is sent with a Content-Type header set to application/json
*/
func TestDefaultContentType(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "application/json", req.Header.Get(HeaderContentType))
		_, err := rw.Write([]byte(`OK`))
		if err != nil {
			return
		}
	}))
	defer server.Close()
	cmp := &DefaultComponent{}
	conf := cmp.Settings()
	tr, _ := cmp.New(context.Background(), conf)
	client := server.Client()
	client.Transport = tr
	req, _ := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(`{"hello": "world"}`))
	resp, _ := client.Do(req)
	assert.Equal(t, 200, resp.StatusCode)
}

/*
Given a default http client with a default Content-Type of application/json
When a request is sent with a Content-Type header set to application/jsonlines
Then the request is sent with a Content-Type header set to application/jsonlines
*/
func TestDefaultContentTypeOverrideable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "application/jsonlines", req.Header.Get(HeaderContentType))
		_, err := rw.Write([]byte(`OK`))
		if err != nil {
			return
		}
	}))
	defer server.Close()
	cmp := &DefaultComponent{}
	conf := cmp.Settings()
	tr, _ := cmp.New(context.Background(), conf)
	client := server.Client()
	client.Transport = tr
	req, _ := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(`{"hello": "world"}`))
	req.Header.Set(HeaderContentType, "application/jsonlines")
	resp, _ := client.Do(req)
	assert.Equal(t, 200, resp.StatusCode)
}
