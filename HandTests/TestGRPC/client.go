package main

import (
	// ...
	"context"
	"fmt"
	"log"

	pb "github.com/Arti9991/shortener/internal/gRPC/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	c := pb.NewShortenerClient(conn)

	// функция, в которой будем отправлять сообщения
	TestShortener(c)
	TestShortenerJSON(c)
	TestShortenerUser(c)
}

func TestShortener(c pb.ShortenerClient) {
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
		log.Fatal(err)
	}
	fmt.Println(resp2.Addres)

	// md2 := metadata.New(map[string]string{"UserID": "123123143124"})
	// ctx2 := metadata.NewOutgoingContext(context.Background(), md2)

	resp3, err := c.GetAddrUser(ctx, &pb.GetAddrUserRequset{}, grpc.Header(&header))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp3.Addreses)

	// if resp.Error != "" {
	// 	fmt.Println(resp.Error)
	// }
	// }
	// // удаляем одного из пользователей
	// resp, err := c.DelUser(context.Background(), &pb.DelUserRequest{
	// 	Email: "serge@example.com",
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if resp.Error != "" {
	// 	fmt.Println(resp.Error)
	// }
}

func TestShortenerJSON(c pb.ShortenerClient) {
	// набор тестовых данных
	// for _, user := range users {
	md := metadata.New(map[string]string{"UserID": "8c537969b84ad4eb0a73e29b3f2a9030"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var header metadata.MD
	// добавляем пользователей
	resp, err := c.PostAddr(ctx, &pb.PostAddrRequset{
		Addres: `{"url":"www.Dlya.ru"}`,
		IsJSON: true,
	},
		grpc.Header(&header))

	if err != nil {
		log.Fatal(err)
	}
	UserID := header.Get("UserID")
	shortAddr := resp.Addres
	fmt.Println(UserID)
	fmt.Println(shortAddr)
}

func TestShortenerUser(c pb.ShortenerClient) {

	//md := metadata.New(map[string]string{"UserID": "8c537969b84ad4eb0a73e29b3f2a9030"})
	md := metadata.New(map[string]string{"UserID": "123456"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var header metadata.MD

	resp3, err := c.GetAddrUser(ctx, &pb.GetAddrUserRequset{}, grpc.Header(&header))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp3.Addreses)
	UserID := header.Get("UserID")
	fmt.Println(UserID)
	// if resp.Error != "" {
	// 	fmt.Println(resp.Error)
	// }
	// }
	// // удаляем одного из пользователей
	// resp, err := c.DelUser(context.Background(), &pb.DelUserRequest{
	// 	Email: "serge@example.com",
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if resp.Error != "" {
	// 	fmt.Println(resp.Error)
	// }
}
