package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func Example() {
	// создаем тестовые струткуры для сервера.
	server := NewServerTest()
	test := httptest.NewServer(server.MainRouter())
	// создаем тестовый клиент с отклюенным redirect
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// подготавливаем POST запрос с URL в теле запроса.
	request, err := http.NewRequest("POST", test.URL, bytes.NewBuffer([]byte("www.smth.com")))
	if err != nil {
		return
	}
	request.Header.Add("Content-Type", "text/plain")
	// выполняем запрос
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Результат первого запроса:")
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	// читаем тело ответа.
	respBody, err := io.ReadAll(resp.Body)
	// убираем базовый URL (оставляем только hash для удобства).
	hash, _ := strings.CutPrefix(string(respBody), "http://example.com/")
	resp.Body.Close()
	// сохраняем user ID из полученной куки.
	cookie := resp.Cookies()[0]

	// делаеми get запрос с укороченным URL для получения исходного URL.
	request2, err := http.NewRequest("GET", test.URL+"/"+hash, nil)
	if err != nil {
		return
	}
	resp2, err := client.Do(request2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Результат второго запроса:")
	fmt.Println(resp2.StatusCode)
	fmt.Println(resp2.Header.Get("Location"))
	resp2.Body.Close()

	// делаем запрос для удаления URL по его укороченной верссии полученной в первом запросе.
	request3, err := http.NewRequest("DELETE", test.URL+"/api/user/urls", bytes.NewBuffer([]byte(`["`+hash+`"]`)))
	request3.Header.Add("Content-Type", "application/json")
	// авторизуем пользователя по куке из первого запроса для удаления.
	request3.AddCookie(cookie)
	resp3, err := client.Do(request3)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Результат третьего запроса:")
	fmt.Println(resp3.StatusCode)
	resp3.Body.Close()

	// подготавливаем JSON с несколькими URL для сокращения.
	var batchBody = `[
							{
								"correlation_id": "ID1",
								"original_url": "www.first.ru"
							},
							{
								"correlation_id": "ID2",
								"original_url": "www.second.ru"
							},
							{
								"correlation_id": "ID3",
								"original_url": "www.third.ru"
							},
							{
								"correlation_id": "ID4",
								"original_url": "www.fourth.ru"
							},
							{
								"correlation_id": "ID5",
								"original_url": "www.fifth.ru"
							}
						]`
	// делаем запрос на сохранение сразу нескольких URL.
	request4, err := http.NewRequest("POST", test.URL+"/api/shorten/batch", bytes.NewBuffer([]byte(batchBody)))
	if err != nil {
		return
	}
	request4.Header.Add("Content-Type", "application/json")
	// сохранение производим от имени пользователя из первого запроса.
	request.AddCookie(cookie)
	resp4, err := client.Do(request4)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Результат четвертого запроса:")
	fmt.Println(resp4.StatusCode)
	fmt.Println(resp4.Header.Get("Content-Type"))
	resp4.Body.Close()
	// Output:
	// Результат первого запроса:
	// 201
	// text/plain
	// Результат второго запроса:
	// 307
	// www.smth.com
	// Результат третьего запроса:
	// 202
	// Результат четвертого запроса:
	// 201
	// application/json
}
