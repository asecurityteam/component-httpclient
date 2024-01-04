package httpclient

import (
	"context"
	"testing"

	"github.com/asecurityteam/settings"
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
