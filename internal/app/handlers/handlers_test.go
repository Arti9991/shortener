package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/Arti9991/shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainPage(t *testing.T) {
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
				answer:      "http://localhost:8080/",
			},
		},
		{
			name:    "Test for error with no body",
			request: "/",
			body:    "",
			want: want{
				statusCode:  400,
				contentType: "text/plain",
				answer:      "http://localhost:8080/",
			},
		},
		{
			name:    "Test for unusual request",
			request: "/sdadefedfsaa",
			body:    "",
			want: want{
				statusCode:  400,
				contentType: "text/plain",
				answer:      "http://localhost:8080/",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(MainPage(&dt))
			h(w, request)

			result := w.Result()
			if result.StatusCode != http.StatusBadRequest {
				assert.Equal(t, test.want.statusCode, result.StatusCode)
				assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))

				userResult, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)

				re := regexp.MustCompile(`http://example.com/\w+`)
				strResult := string(userResult)
				assert.True(t, re.MatchString(strResult))
			} else if result.StatusCode == http.StatusBadRequest {
				assert.Equal(t, test.want.statusCode, result.StatusCode)
			}

		})
	}
}

func TestAllHandle(t *testing.T) {
	type want struct {
		statusCode1  int
		statusCode2  int
		contentType1 string
		contentType2 string
		location     string
	}
	dt := storage.NewData()
	tests := []struct {
		name    string
		request string
		body    string
		want    want
	}{
		{
			name:    "Simple request for code 307",
			request: "/",
			body:    "www.ya.ru",
			want: want{
				statusCode1:  201,
				statusCode2:  307,
				contentType1: "text/plain",
				contentType2: "",
				location:     "www.ya.ru",
			},
		},
		{
			name:    "Non formal request for code 307",
			request: "/",
			body:    "Quentin Tarantino #4sd sd4fr 4d54354",
			want: want{
				statusCode1:  201,
				statusCode2:  307,
				contentType1: "text/plain",
				contentType2: "",
				location:     "Quentin Tarantino #4sd sd4fr 4d54354",
			},
		},
		{
			name:    "Test for error with no body",
			request: "/",
			body:    "",
			want: want{
				statusCode1: 400,
			},
		},
		{
			name:    "Test for unusual request",
			request: "/sdadefedfsaa",
			body:    "",
			want: want{
				statusCode1: 400,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request1 := httptest.NewRequest(http.MethodPost, test.request, strings.NewReader(test.body))
			w1 := httptest.NewRecorder()
			h1 := http.HandlerFunc(MainPage(&dt))
			h1(w1, request1)
			result := w1.Result()

			if result.StatusCode != http.StatusBadRequest {
				assert.Equal(t, test.want.statusCode1, result.StatusCode)
				assert.Equal(t, test.want.contentType1, result.Header.Get("Content-Type"))

				userResult, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)

				re := regexp.MustCompile(`http://example.com/\w+`)
				strResult := string(userResult)
				assert.True(t, re.MatchString(strResult))

				request2 := httptest.NewRequest(http.MethodGet, strResult, nil)
				w2 := httptest.NewRecorder()
				h2 := http.HandlerFunc(GetAddr(&dt))
				h2(w2, request2)
				result2 := w2.Result()

				assert.Equal(t, test.want.statusCode2, result2.StatusCode)
				assert.Equal(t, test.want.contentType2, result2.Header.Get("Content-Type"))
				assert.Equal(t, test.want.location, result2.Header.Get("Location"))

				err = result2.Body.Close()
				require.NoError(t, err)
			} else if result.StatusCode == http.StatusBadRequest {
				assert.Equal(t, test.want.statusCode1, result.StatusCode)
			}
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
				h1 := http.HandlerFunc(MainPage(&dt))
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
			}
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

				fmt.Printf("\n\n\nIter: %d Result: %s\n\n\n", i, result2.Header.Get("Location"))
			}
		})
	}
}
