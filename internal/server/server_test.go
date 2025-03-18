package server

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Arti9991/shortener/internal/config"
	"github.com/Arti9991/shortener/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

// инциализация подменной структуры для тестов с отключенным логгированием
func NewServerTest() *Server {
	// установка сида для случайных чисел
	rand.Seed(uint64(time.Now().UnixNano()))
	var Serv Server
	// инициализация конфигурации
	Serv.Config = config.InitConfTests()
	// инициализация логгера
	err := logger.Initialize("Error")
	if err != nil {
		panic(err)
	}
	Serv.StorInit()

	return &Serv
}

func findValue(res string, wants []string) bool {
	for _, want := range wants {
		if res == want {
			return true
		}
	}
	return false
}

var serv = NewServerTest()
var ts = httptest.NewServer(serv.MainRouter())

func testRequests(t *testing.T, ts *httptest.Server, method,
	path string, body io.Reader) (*http.Response, string) {
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

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	///
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
		{
			name:    "Seven non formal requests for code 307",
			request: "/",
			bodys: []string{
				"Booogie woogie",
				"Gophers cool",
				"Nice look",
				"Etc",
				"Non Formal tests",
				"Errors",
				"Boooo",
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
					"Errors",
					"Boooo",
				},
			},
		},
	}
	for _, test := range tests {
		ident := make([]string, 0)
		locate := make([]string, 0)
		for i := range len(test.bodys) {
			resp, get := testRequests(t, ts, "POST", test.request, strings.NewReader(test.bodys[i]))
			defer resp.Body.Close()
			assert.Equal(t, test.want.statusCode1, resp.StatusCode)
			assert.Equal(t, test.want.contentType1, resp.Header.Get("Content-Type"))
			get, found := strings.CutPrefix(get, "http://example.com")
			assert.True(t, found)
			ident = append(ident, get)
			resp2, _ := testRequests(t, ts, "GET", ident[i], nil)
			defer resp2.Body.Close()

			assert.Equal(t, test.want.statusCode2, resp2.StatusCode)
			assert.Equal(t, test.want.contentType2, resp2.Header.Get("Content-Type"))
			locate = append(locate, resp2.Header.Get("Location"))
		}
		time.Sleep(50 * time.Millisecond)
		for _, loc := range locate {
			assert.True(t, findValue(loc, test.want.locations))
		}

	}
}

func testRequestCompress(t *testing.T, ts *httptest.Server, method,
	path string, body io.Reader) (*http.Response, string) {
	client := &http.Client{}
	request, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Content-Encoding", "gzip")
	request.Header.Add("Accept-Encoding", "gzip")

	resp, err := client.Do(request)
	require.NoError(t, err)
	defer resp.Body.Close()

	zr, err := gzip.NewReader(resp.Body)
	require.NoError(t, err)

	b, err := io.ReadAll(zr)
	require.NoError(t, err)

	ResURL := &struct {
		Result string `json:"result"`
	}{}
	err = json.Unmarshal(b, ResURL)
	require.NoError(t, err)

	return resp, string(ResURL.Result)
}
func TestRouterCompress(t *testing.T) {
	///
	type want struct {
		statusCode1     int
		contentType     string
		contentEncoding string
	}
	tests := []struct {
		name    string
		request string
		body    string
		want    want
	}{
		{
			name:    "Simple request for code 201 with compression for request and responce",
			request: "/api/shorten",
			body:    `{"url":"www.ya.ru"}`,
			want: want{
				statusCode1:     201,
				contentType:     "application/json",
				contentEncoding: "gzip",
			},
		},
		{
			name:    "Long request for code 201 with compression for request and responce",
			request: "/api/shorten",
			body:    `{"url":"booly/boolean/true/means/23452dsaf432drfredt43fkpymejudmnr4pgjrvotmgi"}`,
			want: want{
				statusCode1:     201,
				contentType:     "application/json",
				contentEncoding: "gzip",
			},
		},
	}
	for _, test := range tests {
		//ident := make([]string, 0)

		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(test.body))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		resp, get := testRequestCompress(t, ts, "POST", test.request, buf)
		defer resp.Body.Close()

		assert.Equal(t, test.want.statusCode1, resp.StatusCode)
		assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		assert.Equal(t, test.want.contentEncoding, resp.Header.Get("Content-Encoding"))

		_, found := strings.CutPrefix(get, "http://example.com")
		assert.True(t, found)
	}
}

func testRequestsBench(b *testing.B, ts *httptest.Server, method,
	path string, body io.Reader) (*http.Response, string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	request, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(b, err)

	resp, err := client.Do(request)
	require.NoError(b, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(b, err)

	return resp, string(respBody)
}

func BenchmarkServer(b *testing.B) {
	///
	type want struct {
		statusCode1  int
		statusCode2  int
		contentType1 string
		contentType2 string
		location     string
	}
	tests := []struct {
		name    string
		request string
		body    string
		want    want
	}{
		{
			name:    "Simple request for code 307",
			request: "/",
			body:    "www.smth.ru",
			want: want{
				statusCode1:  201,
				statusCode2:  307,
				contentType1: "text/plain",
				contentType2: "",
				location:     "www.smth.ru",
			},
		},
	}
	b.ResetTimer()
	for _, test := range tests {
		var ident string
		b.Run("POST", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				resp, get := testRequestsBench(b, ts, "POST", test.request, strings.NewReader(test.body))
				b.StopTimer()

				defer resp.Body.Close()
				assert.Equal(b, test.want.statusCode1, resp.StatusCode)
				assert.Equal(b, test.want.contentType1, resp.Header.Get("Content-Type"))
				get, found := strings.CutPrefix(get, "http://example.com")
				assert.True(b, found)
				ident = get
			}
		})
		b.Run("GET", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				resp2, _ := testRequestsBench(b, ts, "GET", ident, nil)
				b.StopTimer()
				defer resp2.Body.Close()
				assert.Equal(b, test.want.statusCode2, resp2.StatusCode)
				assert.Equal(b, test.want.location, resp2.Header.Get("Location"))
			}
		})

	}
}
