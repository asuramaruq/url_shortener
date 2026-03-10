package save

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asuramaruq/url_shortener/internal/http-server/handlers/url/save/mocks"
	"github.com/asuramaruq/url_shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError string
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlSaverMock := mocks.NewURLSaver(t)

			if tc.respError == "" || tc.mockError != "" {
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).Return(int64(1), tc.mockError).Once()
			}

			handler := New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))

			require.NoError(t, err)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			var resp Response

			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
