package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Arti9991/shortener/internal/storage/files"
	"github.com/Arti9991/shortener/internal/storage/inmemory"
	"github.com/Arti9991/shortener/internal/storage/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//для базовых тестов производится генерация моков командой ниже
// mockgen --source=./internal/storage/storage.go --destination=./internal/storage/mocks/mocks_store.go --package=mocks StorFunc

var BaseAdr = "http://example.com"
var Files = files.FilesTest()

func TestPostAddr(t *testing.T) {
	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	m := mocks.NewMockStorFunc(ctrl)

	// задаем режим рабоыт моков (для POST главное отсутствие ошибки)
	m.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil).
		MaxTimes(1)
	hd := NewHandlersData(m, BaseAdr, files.FilesTest())

	type want struct {
		statusCode  int
		contentType string
		answer      string
	}

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
			h := http.HandlerFunc(PostAddr(hd))
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
		})
	}
}

func TestPostAddrJSON(t *testing.T) {
	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	m := mocks.NewMockStorFunc(ctrl)
	// задаем режим рабоыт моков (для POST главное отсутствие ошибки)
	m.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil).
		MaxTimes(2)
	hd := NewHandlersData(m, BaseAdr, files.FilesTest())

	type want struct {
		statusCode  int
		contentType string
		answer      string
	}
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
			h := http.HandlerFunc(PostAddrJSON(hd))
			h(w, request)

			result := w.Result()
			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))

			err := json.NewDecoder(result.Body).Decode(&ResURL)
			require.NoError(t, err)

			strResult := string(ResURL.Result)
			fmt.Println(strResult)

			// res, _ := strings.CutPrefix(strResult, "http://example.com/")
			// assert.Equal(t, test.want.answer, hd.Dt.ShortUrls[res])
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestGet(t *testing.T) {
	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	m := mocks.NewMockStorFunc(ctrl)

	type want struct {
		statusCode int
		answer     string
		err        error
	}
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
				err:        nil,
			},
		},
		{
			name:    "Big request for code 307",
			hash:    "FXFGaseD",
			request: "/FXFGaseD",
			want: want{
				statusCode: 307,
				answer:     "/env/local/path_slide/beta",
				err:        nil,
			},
		},
		{
			name:    "Test for error with no hash",
			hash:    "/",
			request: "/",
			want: want{
				statusCode: 400,
				answer:     "",
				err:        errors.New("no such URL in memory"),
			},
		},
		{
			name:    "Test for error with bad hash",
			hash:    "SAGREVad",
			request: "/SAGREVad",
			want: want{
				statusCode: 400,
				answer:     "",
				err:        errors.New("no such URL in memory"),
			},
		},
	}
	for _, test := range tests {
		// задаем режим рабоыт моков (для GET проверяем полученные файлы)
		m.EXPECT().
			Get(test.hash).
			Return(test.want.answer, test.want.err).
			MaxTimes(1)

		hd := NewHandlersData(m, BaseAdr, files.FilesTest())

		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetAddr(hd))
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
	// для сложных запросов используем подменную структуру с хранением данных в памяти
	hd := NewHandlersData(inmemory.NewData(Files), BaseAdr, files.FilesTest())

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
				h1 := http.HandlerFunc(PostAddr(hd))
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
				h2 := http.HandlerFunc(GetAddr(hd))
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

func TestPostBatch(t *testing.T) {
	// для сложных запросов используем подменную структуру с хранением данных в памяти
	hd := NewHandlersData(inmemory.NewData(Files), BaseAdr, files.FilesTest())

	type want struct {
		statusCode  int
		contentType string
		answers     []string
	}
	tests := []struct {
		name    string
		request string
		income  string
		want    want
	}{
		{
			name:    "Multiple requests in one JSON for code 201",
			request: "/api/shorten/batch",
			income: `[
							{
								"correlation_id": "ID",
								"original_url": "www.ya.ru"
							},
							{
								"correlation_id": "ID",
								"original_url": "www.dlya.ru"
							},
							{
								"correlation_id": "ID",
								"original_url": "www.Nya.ru"
							},
							{
								"correlation_id": "ID",
								"original_url": "www.Qya.ru"
							},
							{
								"correlation_id": "ID",
								"original_url": "www.Mya.ru"
							}
						]`,
			want: want{
				statusCode:  201,
				contentType: "application/json",
				answers: []string{
					"www.ya.ru",
					"www.dlya.ru",
					"www.Nya.ru",
					"www.Qya.ru",
					"www.Mya.ru",
				},
			},
		},
		{
			name:    "One request in JSON for code 201",
			request: "/api/shorten/batch",
			income: `[
							{
								"correlation_id": "ID",
								"original_url": "www.ya.ru"
							}
						]`,
			want: want{
				statusCode:  201,
				contentType: "application/json",
				answers: []string{
					"www.ya.ru",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			ResURL := []struct {
				CorrID   string `json:"correlation_id"`
				ShortURL string `json:"short_url"`
			}{}
			request := httptest.NewRequest(http.MethodPost, test.request, bytes.NewBuffer([]byte(test.income)))
			request.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			h := http.HandlerFunc(PostBatch(hd))
			h(w, request)

			result := w.Result()
			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))

			err := json.NewDecoder(result.Body).Decode(&ResURL)
			require.NoError(t, err)

			for i := range len(ResURL) {

				strResult := string(ResURL[i].ShortURL)

				res, _ := strings.CutPrefix(strResult, "http://example.com/")
				inmem, _ := hd.Dt.Get(res)
				assert.Equal(t, test.want.answers[i], inmem)
			}
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}
