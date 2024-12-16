package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Arti9991/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostAddr(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		answer      string
	}
	dt := storage.NewData()
	tests := []struct {
		name    string
		request string
		body    string
		want    want
	}{
		{
			name:    "Simple request for code 201",
			request: "/",
			body:    "www.ya.ru",
			want: want{
				statusCode:  201,
				contentType: "text/plain",
				answer:      "http://example.com/",
			},
		},
		{
			name:    "Test for error with no body",
			request: "/",
			body:    "",
			want: want{
				statusCode:  400,
				contentType: "",
				answer:      "",
			},
		},
		{
			name:    "Test for unusual request",
			request: "/sdadefedfsaa",
			body:    "",
			want: want{
				statusCode:  400,
				contentType: "",
				answer:      "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(PostAddr(&dt, "http://example.com"))
			h(w, request)

			result := w.Result()
			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))

			userResult, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			strResult := string(userResult)
			bl := strings.Contains(strResult, test.want.answer)
			assert.True(t, bl)
			res, _ := strings.CutPrefix(strResult, "http://example.com/")
			assert.Equal(t, test.body, dt.ShortUrls[res])
		})
	}
}

func TestPostAddrJSON(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		answer      string
	}
	dt := storage.NewData()
	tests := []struct {
		name    string
		request string
		income  string
		want    want
	}{
		{
			name:    "Simple request for code 201",
			request: "/api/shorten",
			income:  `{"url":"www.ya.ru"}`,
			want: want{
				statusCode:  201,
				contentType: "application/json",
				answer:      "www.ya.ru",
			},
		},
		{
			name:    "Long request for code 201",
			request: "/api/shorten",
			income:  `{"url":"passpot/idcheck/name/definition/correct"}`,
			want: want{
				statusCode:  201,
				contentType: "application/json",
				answer:      "passpot/idcheck/name/definition/correct",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			ResURL := &struct {
				Result string `json:"result"`
			}{}

			request := httptest.NewRequest(http.MethodPost, test.request, bytes.NewBuffer([]byte(test.income)))
			request.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			h := http.HandlerFunc(PostAddrJSON(&dt, "http://example.com"))
			h(w, request)

			result := w.Result()
			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))

			err := json.NewDecoder(result.Body).Decode(&ResURL)
			require.NoError(t, err)

			//fmt.Println(ResURL.Result)
			strResult := string(ResURL.Result)

			res, _ := strings.CutPrefix(strResult, "http://example.com/")
			assert.Equal(t, test.want.answer, dt.ShortUrls[res])
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestGet(t *testing.T) {
	type want struct {
		statusCode int
		answer     string
	}
	dt := storage.NewData()
	tests := []struct {
		name    string
		hash    string
		request string
		want    want
	}{
		{
			name:    "Simple request for code 307",
			hash:    "DxDfgvDa",
			request: "/DxDfgvDa",
			want: want{
				statusCode: 307,
				answer:     "ya.ru",
			},
		},
		{
			name:    "Big request for code 307",
			hash:    "FXFGaseD",
			request: "/FXFGaseD",
			want: want{
				statusCode: 307,
				answer:     "/env/local/path_slide/beta",
			},
		},
		{
			name:    "Test for error with no hash",
			hash:    "AMFhvnth",
			request: "/",
			want: want{
				statusCode: 400,
				answer:     "",
			},
		},
		{
			name:    "Test for error with bad hash",
			hash:    "DxDfgvDa",
			request: "/SAGREVad",
			want: want{
				statusCode: 400,
				answer:     "",
			},
		},
	}
	for _, test := range tests {
		dt.AddValue(test.hash, test.want.answer)
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetAddr(&dt))
			h(w, request)
			result := w.Result()
			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Equal(t, test.want.answer, result.Header.Get("Location"))

			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestMultuplTasks(t *testing.T) {
	type want struct {
		statusCode1  int
		statusCode2  int
		contentType1 string
		contentType2 string
		locations    []string
	}
	dt := storage.NewData()
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
			name:    "Three requests for code 307",
			request: "/",
			bodys: []string{
				"www.ya.ru",
				"www.vk,com",
				"www.rbk.ru",
			},
			want: want{
				statusCode1:  201,
				statusCode2:  307,
				contentType1: "text/plain",
				contentType2: "",
				locations: []string{
					"www.ya.ru",
					"www.vk,com",
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
		t.Run(test.name, func(t *testing.T) {
			strResults := make([]string, 0)
			for _, body := range test.bodys {
				request1 := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(body))
				w1 := httptest.NewRecorder()
				h1 := http.HandlerFunc(PostAddr(&dt, "http://example.com"))
				h1(w1, request1)
				result := w1.Result()

				assert.Equal(t, test.want.statusCode1, result.StatusCode)
				assert.Equal(t, test.want.contentType1, result.Header.Get("Content-Type"))

				userResult, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)

				re := regexp.MustCompile(`http://example.com/\w+`)
				strResult := string(userResult)
				assert.True(t, re.MatchString(strResult))
				strResults = append(strResults, strResult)
				time.Sleep(10 * time.Millisecond)
			}

			time.Sleep(100 * time.Millisecond)
			for i, loc := range test.want.locations {
				request2 := httptest.NewRequest(http.MethodGet, strResults[i], nil)
				w2 := httptest.NewRecorder()
				h2 := http.HandlerFunc(GetAddr(&dt))
				h2(w2, request2)
				result2 := w2.Result()

				assert.Equal(t, test.want.statusCode2, result2.StatusCode)
				assert.Equal(t, test.want.contentType2, result2.Header.Get("Content-Type"))
				assert.Equal(t, loc, result2.Header.Get("Location"))

				err := result2.Body.Close()
				require.NoError(t, err)

			}
		})
	}
}
