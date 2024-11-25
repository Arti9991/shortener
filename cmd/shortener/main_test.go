// main_test.go
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Arti9991/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequests(t *testing.T, ts *httptest.Server, method,
	path string, body io.Reader) (*http.Response, string, string) {
	//req, err := http.NewRequest(method, ts.URL+path, nil)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	request, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := client.Do(request)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	reqURL := request.URL.String()
	return resp, string(respBody), reqURL
}

func TestRouter(t *testing.T) {
	data := storage.NewData()
	ts := httptest.NewServer(MainRouter(data))
	defer ts.Close()

	type want struct {
		statusCode1  int
		statusCode2  int
		contentType1 string
		contentType2 string
		locations    []string
	}
	tests := []struct {
		name    string
		request string
		bodys   []string
		want    want
	}{
		{
			name:    "Simple request for code 307",
			request: "/",
			bodys:   []string{"www.eto.ne.ya.ru"},
			want: want{
				statusCode1:  201,
				statusCode2:  307,
				contentType1: "text/plain",
				contentType2: "",
				locations:    []string{"www.eto.ne.ya.ru"},
			},
		},
		{
			name:    "Non formal request for code 307",
			request: "/",
			bodys:   []string{"Quentin Tarantino #4sd sd4fr 4d54354"},
			want: want{
				statusCode1:  201,
				statusCode2:  307,
				contentType1: "text/plain",
				contentType2: "",
				locations:    []string{"Quentin Tarantino #4sd sd4fr 4d54354"},
			},
		},
		{
			name:    "Three requests for code 307",
			request: "/",
			bodys: []string{
				"www.ya.ru",
				"www.vk.com",
				"www.rbk.ru",
			},
			want: want{
				statusCode1:  201,
				statusCode2:  307,
				contentType1: "text/plain",
				contentType2: "",
				locations: []string{
					"www.ya.ru",
					"www.vk.com",
					"www.rbk.ru",
				},
			},
		},
		{
			name:    "Five non formal requests for code 307",
			request: "/",
			bodys: []string{
				"Booogie woogie",
				"Gophers cool",
				"Nice look",
				"Etc",
				"Non Formal tests",
			},
			want: want{
				statusCode1:  201,
				statusCode2:  307,
				contentType1: "text/plain",
				contentType2: "",
				locations: []string{
					"Booogie woogie",
					"Gophers cool",
					"Nice look",
					"Etc",
					"Non Formal tests",
				},
			},
		},
	}
	for _, test := range tests {
		ident := make([]string, 0)
		for i := range len(test.bodys) {
			resp, get, reqURL := testRequests(t, ts, "POST", test.request, strings.NewReader(test.bodys[i]))
			assert.Equal(t, test.want.statusCode1, resp.StatusCode)
			assert.Equal(t, test.want.contentType1, resp.Header.Get("Content-Type"))
			bl := strings.Contains(get, reqURL)
			assert.True(t, bl)
			// assert.Equal(t, v.status, resp.StatusCode)
			// assert.Equal(t, v.want, get)
			get, found := strings.CutPrefix(get, reqURL)
			assert.True(t, found)
			ident = append(ident, "/"+get)
			// }
			// for i := range len(test.want.locations) {
			fmt.Println(ident)

			resp2, _, _ := testRequests(t, ts, "GET", ident[i], nil)
			// fmt.Printf("\n")
			// fmt.Println(OutUrl)
			// fmt.Println(reqURL)
			// fmt.Printf("\n")

			assert.Equal(t, test.want.statusCode2, resp2.StatusCode)
			assert.Equal(t, test.want.contentType2, resp2.Header.Get("Content-Type"))
			assert.Equal(t, test.want.locations[i], resp2.Header.Get("Location"))
		}
	}
}
