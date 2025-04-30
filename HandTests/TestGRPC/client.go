package main

import (
	// ...
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"regexp"
	"strconv"

	pb "github.com/Arti9991/shortener/internal/gRPC/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {

	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	c := pb.NewShortenerClient(conn)

	// функция, в которой будем отправлять сообщения
	BaseTestShortener(c)
	// TestShortenerJSON(c)
	TestShortenerbatch(c)
	TestShortenerUser(c)
	PingTest(c)
	TestGetStats(c)
}

// BaseTestShortener функция для простейших POST и GET запросов
func BaseTestShortener(c pb.ShortenerClient) {
	// набор тестовых данных
	// for _, user := range users {
	md := metadata.New(map[string]string{"UserID": "8c537969b84ad4eb0a73e29b3f2a9030"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var header metadata.MD
	// добавляем пользователей
	resp, err := c.PostAddr(ctx, &pb.PostAddrRequset{
		Addres: "www.ya.ru",
	},
		grpc.Header(&header))

	if err != nil {
		log.Fatal(err)
	}
	UserID := header.Get("UserID")
	shortAddr := resp.Addres
	fmt.Println(UserID)
	fmt.Println(shortAddr)

	resp2, err := c.GetAddr(ctx, &pb.GetAddrRequset{
		ShortAddr: shortAddr,
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp2.Addres)
	}

	// md2 := metadata.New(map[string]string{"UserID": "123123143124"})
	// ctx2 := metadata.NewOutgoingContext(context.Background(), md2)

	resp3, err := c.GetAddrUser(ctx, &pb.GetAddrUserRequset{}, grpc.Header(&header))
	if err != nil {
		log.Fatal(err)
	}
	for _, val := range resp3.UserURLs {
		fmt.Printf("Orig_URL: %s\t Short_URL: %s\n", val.OrigURL, val.ShortURL)
	}
	re := regexp.MustCompile(`^.*/`)
	shortIdent := []string{re.ReplaceAllString(shortAddr, "")}

	_, err = c.DeleteAddr(ctx, &pb.DeleteAddrRequest{Idents: shortIdent}, grpc.Header(&header))
	if err != nil {
		log.Fatal(err)
	}
}

// func TestShortenerJSON(c pb.ShortenerClient) {
// 	// набор тестовых данных
// 	// for _, user := range users {
// 	md := metadata.New(map[string]string{"UserID": "8c537969b84ad4eb0a73e29b3f2a9030"})
// 	ctx := metadata.NewOutgoingContext(context.Background(), md)

// 	var header metadata.MD
// 	// добавляем пользователей
// 	resp, err := c.PostAddr(ctx, &pb.PostAddrRequset{
// 		Addres: `{"url":"www.Dlya.ru"}`,
// 		IsJSON: true,
// 	},
// 		grpc.Header(&header))

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	UserID := header.Get("UserID")
// 	shortAddr := resp.Addres
// 	fmt.Println(UserID)
// 	fmt.Println(shortAddr)
// }

// TestShortenerUser функция для получения информации о пользователе
func TestShortenerUser(c pb.ShortenerClient) {

	//md := metadata.New(map[string]string{"UserID": "8c537969b84ad4eb0a73e29b3f2a9030"})
	md := metadata.New(map[string]string{"UserID": "123"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var header metadata.MD

	resp3, err := c.GetAddrUser(ctx, &pb.GetAddrUserRequset{}, grpc.Header(&header))
	if err != nil {
		log.Fatal(err)
	}

	for _, val := range resp3.UserURLs {
		fmt.Printf("Orig_URL: %s\t Short_URL: %s\n", val.OrigURL, val.ShortURL)
	}
	UserID := header.Get("UserID")
	if len(UserID) > 0 {
		fmt.Printf("New UserID is: %s\n", UserID[0])
	}
}

// TestShortenerbatch функция для Batch POST запроса
func TestShortenerbatch(c pb.ShortenerClient) {
	md := metadata.New(map[string]string{"UserID": "97eb08b6c9e4edf594a43793e825a4b7"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	//var data []byte
	n := 4
	URLs := make([]*pb.BatchURL, n)

	for i := 0; i < n; i++ {
		var Save pb.BatchURL
		str := "unique_URL_" + strconv.Itoa(i) + "_" + strconv.Itoa(rand.IntN(100000)) + ".com"
		Save.URL = str
		Save.CorrID = "ID_" + strconv.Itoa(i)
		URLs[i] = &Save
	}

	//fmt.Println(URLs)

	// data, err := json.Marshal(URLs)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// strDt := `[
	// 	{
	// 		"correlation_id": "ID1",
	// 		"original_url": "www.unque1.ru"
	// 	},
	// 	{
	// 		"correlation_id": "ID2",
	// 		"original_url": "www.unque2.ru"
	// 	},
	// 	{
	// 		"correlation_id": "ID3",
	// 		"original_url": "www.unque3.ru"
	// 	},
	// 	{
	// 		"correlation_id": "ID4",
	// 		"original_url": "www.unque4.ru"
	// 	},
	// 	{
	// 		"correlation_id": "ID5",
	// 		"original_url": "www.unque5.ru"
	// 	}
	// ]`
	// data := []byte(strDt)

	var header metadata.MD
	// добавляем пользователей
	resp, err := c.PostBatch(ctx, &pb.PostBatchRequset{
		BatchURL: URLs,
	},
		grpc.Header(&header))

	if err != nil {
		log.Fatal(err)
	}
	//UserID := header.Get("UserID")

	// var shortAddrs models.OutBuff
	// err = json.Unmarshal(shortAddrsBt, &shortAddrs)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	for _, val := range resp.BatchURL {
		fmt.Printf("Corr_id: %s\t Short_URL: %s\n", val.CorrID, val.URL)
	}
	//fmt.Println(shortAddrs)
	//fmt.Println(shortAddrs)
	if len(resp.BatchURL) > 0 {
		resp2, err := c.GetAddr(ctx, &pb.GetAddrRequset{
			ShortAddr: resp.BatchURL[0].URL,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp2.Addres)
	}
}

// PingTest функция для запроса Ping
func PingTest(c pb.ShortenerClient) {
	md := metadata.New(map[string]string{"UserID": "8c537969b84ad4eb0a73e29b3f2a9030"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := c.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("All fine with DB")
}

// TestGetStats функция для получения статистики по серверу
func TestGetStats(c pb.ShortenerClient) {
	md := metadata.New(map[string]string{"X-Real-IP": "127.0.0.1"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := c.GetStats(ctx, &pb.GetStatsRequest{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Num_users: %d\t Num_URLs: %d\n", resp.NumUsers, resp.NumURLs)
}
