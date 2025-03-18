package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Arti9991/shortener/internal/models"
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

var UserID = "125"
var DeleteChan = make(chan models.DeleteURL)

func TestPostAddr(t *testing.T) {
	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	m := mocks.NewMockStorFunc(ctrl)

	// задаем режим работы моков (для POST главное отсутствие ошибки)
	m.EXPECT().
		Save(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		MaxTimes(1)
	hd := NewHandlersData(m, BaseAdr, files.FilesTest(), DeleteChan)

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

			ctx := context.WithValue(request.Context(), models.CtxKey, models.UserInfo{UserID: UserID})
			request = request.WithContext(ctx)

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
		Save(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		MaxTimes(2)
	hd := NewHandlersData(m, BaseAdr, files.FilesTest(), DeleteChan)

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

			ctx := context.WithValue(request.Context(), models.CtxKey, models.UserInfo{UserID: UserID})
			request = request.WithContext(ctx)

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
				err:        models.ErrorNoURL,
			},
		},
		{
			name:    "Test for error with bad hash",
			hash:    "SAGREVad",
			request: "/SAGREVad",
			want: want{
				statusCode: 400,
				answer:     "",
				err:        models.ErrorNoURL,
			},
		},
	}
	for _, test := range tests {
		// задаем режим рабоыт моков (для GET проверяем полученные файлы)
		m.EXPECT().
			Get(test.hash).
			Return(test.want.answer, test.want.err).
			MaxTimes(1)

		hd := NewHandlersData(m, BaseAdr, files.FilesTest(), DeleteChan)

		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.request, nil)

			ctx := context.WithValue(request.Context(), models.CtxKey, models.UserInfo{UserID: UserID})
			request = request.WithContext(ctx)

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
	hd := NewHandlersData(inmemory.NewData(), BaseAdr, files.FilesTest(), DeleteChan)

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

				ctx1 := context.WithValue(request1.Context(), models.CtxKey, models.UserInfo{UserID: UserID})
				request1 = request1.WithContext(ctx1)

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

				ctx2 := context.WithValue(request2.Context(), models.CtxKey, models.UserInfo{UserID: UserID})
				request2 = request2.WithContext(ctx2)

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
	hd := NewHandlersData(inmemory.NewData(), BaseAdr, files.FilesTest(), DeleteChan)

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

			ctx := context.WithValue(request.Context(), models.CtxKey, models.UserInfo{UserID: UserID})
			request = request.WithContext(ctx)

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

func TestGetUser(t *testing.T) {
	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	m := mocks.NewMockStorFunc(ctrl)

	type want struct {
		statusCode int
		answer     models.UserBuff
		err        error
	}
	tests := []struct {
		name     string
		userID   string
		register bool
		request  string
		want     want
	}{
		{
			name:     "Simple request for code 200",
			userID:   "125",
			register: true,
			request:  "/api/user/urls",
			want: want{
				statusCode: 200,
				answer: []models.UserURL{
					{OrigURL: "www.ya.ru"},
				},
				err: nil,
			},
		},
		{
			name:     "Long request for code 200",
			userID:   "125",
			register: true,
			request:  "/api/user/urls",
			want: want{
				statusCode: 200,
				answer: []models.UserURL{
					{OrigURL: "/env/local/path_slide/beta"},
				},
				err: nil,
			},
		},
		{
			name:     "Many user URLs request for code 200",
			userID:   "125",
			register: true,
			request:  "/api/user/urls",
			want: want{
				statusCode: 200,
				answer: []models.UserURL{
					{OrigURL: "www.ya.ru"},
					{OrigURL: "/env/local/path_slide/beta"},
					{OrigURL: "empty/clown"},
					{OrigURL: "easy"},
					{OrigURL: "/env/local/path_slide/beta/miultiple/twice/long"},
				},
				err: nil,
			},
		},

		{
			name:     "Test for error with good UserID and no URLs",
			userID:   "150",
			register: true,
			request:  "/api/user/urls",
			want: want{
				statusCode: 204,
				answer:     nil,
				err:        models.ErrorNoUserURL,
			},
		},
		{
			name:     "Test for error with bad userID",
			userID:   "150",
			register: false,
			request:  "/api/user/urls",
			want: want{
				statusCode: 204,
				answer:     nil,
				err:        models.ErrorNoUserURL,
			},
		},
	}
	for _, test := range tests {
		// задаем режим рабоыт моков (для GET проверяем полученные файлы)
		m.EXPECT().
			GetUser(test.userID, BaseAdr).
			Return(test.want.answer, test.want.err).
			MaxTimes(5)

		hd := NewHandlersData(m, BaseAdr, files.FilesTest(), DeleteChan)

		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.request, nil)

			ctx := context.WithValue(request.Context(), models.CtxKey, models.UserInfo{
				UserID:   test.userID,
				Register: test.register,
			})
			request = request.WithContext(ctx)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetAddrUser(hd))
			h(w, request)
			result := w.Result()
			assert.Equal(t, test.want.statusCode, result.StatusCode)
			//assert.Equal(t, test.want.answer, result.Header.Get("Location"))

			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestDelete(t *testing.T) {
	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	m := mocks.NewMockStorFunc(ctrl)

	type want struct {
		statusCode int
		err        error
	}
	tests := []struct {
		name     string
		hashes   string
		userID   string
		register bool
		request  string
		want     want
	}{
		{
			name:     "Simple request for code 202",
			hashes:   `["wmmROXSv"]`,
			userID:   "125",
			register: true,
			request:  "/api/user/urls",
			want: want{
				statusCode: 202,
				err:        nil,
			},
		},
		{
			name:     "Many user URLs request for code 202",
			hashes:   `["wmmROXSv","QGtxHvUY","xfZRudbp", "oezaHfOQ", "WsStBGYJ"]`,
			userID:   "125",
			register: true,
			request:  "/api/user/urls",
			want: want{
				statusCode: 202,
				err:        nil,
			},
		},
	}
	for _, test := range tests {
		// задаем режим рабоыт моков (для GET проверяем полученные данные)
		m.EXPECT().
			Delete(gomock.Any(), test.userID).
			Return(test.want.err).
			MaxTimes(5)

		hd := NewHandlersData(m, BaseAdr, files.FilesTest(), DeleteChan)

		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodDelete, test.request, bytes.NewBuffer([]byte(test.hashes)))

			ctx := context.WithValue(request.Context(), models.CtxKey, models.UserInfo{
				UserID:   test.userID,
				Register: test.register,
			})
			request = request.WithContext(ctx)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(DeleteAddr(hd))
			h(w, request)
			result := w.Result()
			assert.Equal(t, test.want.statusCode, result.StatusCode)
			//assert.Equal(t, test.want.answer, result.Header.Get("Location"))

			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
}

func BenchmarkHandlers(b *testing.B) {
	// создаём контроллер
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// создаём объект-заглушку для POST
	m := mocks.NewMockStorFunc(ctrl)

	// задаем режим работы моков для POST
	m.EXPECT().
		Save(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		MinTimes(1)

	// задаем режим работы моков для POST BATCH
	m.EXPECT().
		SaveTx(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		MinTimes(1)

	// задаем режим рабоыт моков для GET
	m.EXPECT().
		Get(gomock.Any()).
		Return("www.ya.ru", nil).
		MinTimes(1)

	// задаем режим рабоыт моков для GET
	// m.EXPECT().
	// 	GetUser(gomock.Any(), gomock.Any()).
	// 	Return(nil, nil).
	// 	MinTimes(1)

	hd := NewHandlersData(m, BaseAdr, files.FilesTest(), DeleteChan)

	// создаём объект-заглушку для GET
	//mG := mocks.NewMockStorFunc(ctrl)

	// // задаем режим рабоыт моков (для GET проверяем полученные файлы)
	// mG.EXPECT().
	// 	Get(test.hash).
	// 	Return(test.want.answer, test.want.err).
	// 	MaxTimes(1)
	// hdG := NewHandlersData(mP, BaseAdr, files.FilesTest(), DeleteChan)

	type want struct {
		statusCodePost    int
		statusCodeGet     int
		statusCodeGetUser int
		contentType       string
		contentTypeJSON   string
		answer            string
	}

	tests := struct {
		name             string
		requestPOST      string
		requestPOSTJSON  string
		requestPOSTBatch string
		requestGET       string
		hash             string
		body             string
		bodyJSON         string
		bodyBatch        string
		want             want
	}{

		name:             "Simple request for code 201",
		requestPOST:      "/",
		requestPOSTJSON:  "/api/shorten",
		requestPOSTBatch: "/api/shorten/batch",
		requestGET:       "/DxDfgvDa",
		hash:             "DxDfgvDa",
		body:             "www.ya.ru",
		bodyJSON:         `{"url":"www.ya.ru"}`,
		bodyBatch: `[
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
			statusCodePost:    201,
			statusCodeGet:     307,
			statusCodeGetUser: 204,
			contentType:       "text/plain",
			contentTypeJSON:   "application/json",
			answer:            "http://example.com/",
		},
	}

	b.Run("POST", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			request := httptest.NewRequest(http.MethodPost, tests.requestPOST, strings.NewReader(tests.body))

			ctx := context.WithValue(request.Context(), models.CtxKey, models.UserInfo{UserID: UserID})
			request = request.WithContext(ctx)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(PostAddr(hd))

			b.StartTimer()
			h(w, request)
			b.StopTimer()

			result := w.Result()
			assert.Equal(b, tests.want.statusCodePost, result.StatusCode)
			assert.Equal(b, tests.want.contentType, result.Header.Get("Content-Type"))

			userResult, err := io.ReadAll(result.Body)
			require.NoError(b, err)
			err = result.Body.Close()
			require.NoError(b, err)

			strResult := string(userResult)
			bl := strings.Contains(strResult, tests.want.answer)
			assert.True(b, bl)
		}
	})

	b.Run("POST JSON", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			request := httptest.NewRequest(http.MethodPost, tests.requestPOSTJSON, bytes.NewBuffer([]byte(tests.bodyJSON)))

			ctx := context.WithValue(request.Context(), models.CtxKey, models.UserInfo{UserID: UserID})
			request = request.WithContext(ctx)

			request.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			h := http.HandlerFunc(PostAddrJSON(hd))

			b.StartTimer()
			h(w, request)
			b.StopTimer()

			result := w.Result()
			assert.Equal(b, tests.want.statusCodePost, result.StatusCode)
			assert.Equal(b, tests.want.contentTypeJSON, result.Header.Get("Content-Type"))

			//strResult := string(ResURL.Result)

			// res, _ := strings.CutPrefix(strResult, "http://example.com/")
			// assert.Equal(t, test.want.answer, hd.Dt.ShortUrls[res])
			err := result.Body.Close()
			require.NoError(b, err)
		}
	})

	b.Run("POST JSON Batch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			request := httptest.NewRequest(http.MethodPost, tests.requestPOSTBatch, bytes.NewBuffer([]byte(tests.bodyBatch)))
			request.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			h := http.HandlerFunc(PostBatch(hd))

			ctx := context.WithValue(request.Context(), models.CtxKey, models.UserInfo{UserID: UserID})
			request = request.WithContext(ctx)

			b.StartTimer()
			h(w, request)
			b.StopTimer()

			result := w.Result()
			assert.Equal(b, tests.want.statusCodePost, result.StatusCode)
			assert.Equal(b, tests.want.contentTypeJSON, result.Header.Get("Content-Type"))

			err := result.Body.Close()
			require.NoError(b, err)
		}
	})

	b.Run("GET", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			request := httptest.NewRequest(http.MethodGet, tests.requestGET, nil)

			ctx := context.WithValue(request.Context(), models.CtxKey, models.UserInfo{UserID: UserID})
			request = request.WithContext(ctx)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetAddr(hd))

			b.StartTimer()
			h(w, request)
			b.StopTimer()

			result := w.Result()
			assert.Equal(b, tests.want.statusCodeGet, result.StatusCode)
			assert.Equal(b, tests.body, result.Header.Get("Location"))

			err := result.Body.Close()
			require.NoError(b, err)
		}
	})
	// b.Run("GET USER", func(b *testing.B) {
	// 	for i := 0; i < b.N; i++ {
	// 		request := httptest.NewRequest(http.MethodGet, tests.requestGET, nil)

	// 		ctx := context.WithValue(request.Context(), models.CtxKey, models.UserInfo{UserID: UserID})
	// 		request = request.WithContext(ctx)

	// 		w := httptest.NewRecorder()
	// 		h := http.HandlerFunc(GetAddrUser(hd))

	// 		b.StartTimer()
	// 		h(w, request)
	// 		b.StopTimer()

	// 		result := w.Result()
	// 		assert.Equal(b, tests.want.statusCodeGetUser, result.StatusCode)

	// 		err := result.Body.Close()
	// 		require.NoError(b, err)
	// 	}
	// })
}
