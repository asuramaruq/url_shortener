package tests

import (
	"net/url"
	"testing"

	"github.com/asuramaruq/url_shortener/internal/http-server/handlers/url/save"
	"github.com/asuramaruq/url_shortener/internal/lib/random"

	"github.com/brianvoe/gofakeit"
	"github.com/gavv/httpexpect"
)

const (
	host = "localhost:8082"
)

func TestURLShortener(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.New(t, u.String())

	e.POST("/url").WithJSON(save.Request{
		URL:   gofakeit.URL(),
		Alias: random.NewRandomString(10),
	}).
		WithBasicAuth("admin", "admin").
		Expect().
		Status(200).
		JSON().
		Object().
		ContainsKey("alias")
}
