package test

import (
	"github.com/kobshobe/agin"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTokenLife(t *testing.T) {
	app := agin.WxApp{
		JwtLife: 1,
	}
	app.Init()
	app.SetJwtLife(1, false)
	token, _ := app.NewJwtToken("testopenid", "1")
	openid, version, err := app.GetTokenInfo(token)
	if err != nil {
		t.Error(err)
	}
	require.Equal(t, "testopenid", openid)
	require.Equal(t, "1", version)
	time.Sleep(time.Second * 2)
	openid, version, err = app.GetTokenInfo(token)
	require.Error(t, err)

	app.SetJwtLife(1, true)
	token, _ = app.NewJwtToken("testopenid", "1")
	openid, version, err = app.GetTokenInfo(token)
	if err != nil {
		t.Error(err)
	}
	require.Equal(t, "testopenid", openid)
	require.Equal(t, "1", version)
	time.Sleep(time.Second * 2)
	openid, version, err = app.GetTokenInfo(token)
	require.NoError(t, err)
}
