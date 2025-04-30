package protoServer

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/Arti9991/shortener/internal/config"
	"github.com/Arti9991/shortener/internal/gRPC/proto"
	pb "github.com/Arti9991/shortener/internal/gRPC/proto"
	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// NewServerTest инциализация подменной структуры для тестов с отключенным логгированием
func InitServerTest() (*ProtoServer, error) {

	var ProtoServ ProtoServer
	var err error
	ctx := context.Background()

	// установка сида для случайных чисел
	rand.Seed(uint64(time.Now().UnixNano()))

	// инициализация конфигурации
	ProtoServ.Config = config.InitConfTests()
	// инициализация логгера
	err = logger.Initialize(ProtoServ.Config.LoggLevel)
	if err != nil {
		return nil, err
	}
	// wait group для ожидания завершения горутин
	// у хэндлеров
	var wgWgHandler sync.WaitGroup
	// инциализация хранилища с нужным интерфейсом
	ProtoServ.StorInit(ctx, &wgWgHandler, ProtoServ.Config.TrustedNet)

	return &ProtoServ, nil
}

func RunGRPCServerTest() error {
	// контекст для ожидания системного сигнала на завершение работы
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// канал для сообщения о Shutdown
	shutCh := make(chan struct{})
	// Wait Group для ожидания завершения горутины удаления
	var WgStor sync.WaitGroup

	serv, err := InitServerTest()
	if err != nil {
		return err
	}

	// определяем адрес для сервера
	listen, err := net.Listen("tcp", serv.Config.HostAdr)
	if err != nil {
		return err
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer(grpc.UnaryInterceptor(atuhInterceptor))
	// регистрируем сервис

	proto.RegisterShortenerServer(s, serv)

	logger.Log.Info("New server initialyzed!",
		zap.String("Server addres:", serv.Config.HostAdr),
		zap.String("Base addres:", serv.Config.BaseAdr),
	)

	// запуск горутины (описана в initStor.go).
	WgStor.Add(1)
	server.RunDeleteStor(serv.Hd, &WgStor)

	// запуск функции ожидающей сигнала о завершении
	// (описана в initStor.go).
	RunWaitShutDown(ctx, s, shutCh)

	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		return err
	}

	// ожидание сообщения о Shutdown
	<-shutCh
	// ожидания закрытия горутины у хэндлера
	serv.Hd.Wg.Wait()
	// закрытие канала отправки URL под удаление
	close(serv.Hd.OutDelCh)
	// ожидание остановки горутины с функцией удаления
	WgStor.Wait()
	// закртытие соединения с базой данных
	err = serv.Hd.Dt.CloseDB()
	if err != nil {
		logger.Log.Info("Error in database close", zap.Error(err))
	}
	logger.Log.Info("Server shutted down!")
	return nil
}

// TestServer интеграционный тест с запуском сервера и запросами к нему
func TestServer(t *testing.T) {
	// запуск горутины с сервером
	go func() {
		err := RunGRPCServerTest()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// ожидание запуска сервера
	time.Sleep(1000 * time.Millisecond)
	// устеновка переменной адреса хоста
	host := "localhost:8080"

	type want struct {
		RespBody      string
		RespBodyBatch []*pb.BatchURL
		NumUrls       int64
		NumUsers      int64
		Err           error
	}
	tests := []struct {
		name           string
		request        string
		UserID         string
		UserIP         string
		BodyPostGet    string
		BodyDeleteAddr []string
		BodyBatch      []*pb.BatchURL
		want           want
	}{
		{
			name:        "Simple POST request with UserID",
			request:     "ordinary post",
			UserID:      "97eb08b6c9e4edf594a43793e825a4b7",
			BodyPostGet: "www.test1.ru",
			want: want{
				RespBody: "http://example.com",
				Err:      nil,
			},
		},
		{
			name:        "Simple POST and GET request with UserID",
			request:     "post amd get",
			UserID:      "97eb08b6c9e4edf594a43793e825a4b7",
			BodyPostGet: "www.test2.ru",
			want: want{
				RespBody: "http://example.com",
				Err:      nil,
			},
		},
		{
			name:    "Batch POST and GET User URLs request with UserID",
			request: "batch post",
			UserID:  "8c537969b84ad4eb0a73e29b3f2a9030",
			BodyBatch: []*pb.BatchURL{
				{
					CorrID: "ID3",
					URL:    "www.test3.ru",
				},
				{
					CorrID: "ID4",
					URL:    "www.test4.ru",
				},
				{
					CorrID: "ID5",
					URL:    "www.test5.ru",
				},
			},
			want: want{
				RespBodyBatch: []*pb.BatchURL{
					{
						CorrID: "ID3",
						URL:    "http://example.com",
					},
					{
						CorrID: "ID4",
						URL:    "http://example.com",
					},
					{
						CorrID: "ID5",
						URL:    "http://example.com",
					},
				},
				Err: nil,
			},
		},
		{
			name:        "GET request with bad user",
			request:     "bad user",
			UserID:      "123",
			BodyPostGet: "www.test6.ru",
			want: want{
				RespBody: "http://example.com",
				Err:      nil,
			},
		},
		{
			name:    "Batch POST, DELETE and GET deleted URL",
			request: "delete addr",
			UserID:  "8c537969b84ad4eb0a73e29b3f2a9030",
			BodyBatch: []*pb.BatchURL{
				{
					CorrID: "ID7",
					URL:    "www.test7.ru",
				},
				{
					CorrID: "ID7",
					URL:    "www.test8.ru",
				},
				{
					CorrID: "ID9",
					URL:    "www.test9.ru",
				},
			},
			want: want{
				Err: nil,
			},
		},
		{
			name:    "GET stats",
			request: "get stats",
			UserID:  "8c537969b84ad4eb0a73e29b3f2a9030",
			UserIP:  "127.0.0.1",
			want: want{
				NumUrls:  6,
				NumUsers: 3,
				Err:      nil,
			},
		},
		// {
		// 	name:    "Non formal request for code 307",
		// 	request: "/",
		// 	bodys:   []string{"Quentin Tarantino #4sd sd4fr 4d54354"},
		// 	want: want{
		// 		statusCode1:  201,
		// 		statusCode2:  307,
		// 		contentType1: "text/plain",
		// 		contentType2: "",
		// 		locations:    []string{"Quentin Tarantino #4sd sd4fr 4d54354"},
		// 	},
		// },
	}
	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	c := pb.NewShortenerClient(conn)

	for _, test := range tests {
		switch test.request {
		case "ordinary post":
			// тест на одиночный запрос POST
			fmt.Println(test.name)
			md := metadata.New(map[string]string{"UserID": test.UserID})
			ctx := metadata.NewOutgoingContext(context.Background(), md)

			var header metadata.MD
			// добавляем пользователей
			resp, err := c.PostAddr(ctx, &pb.PostAddrRequset{
				Addres: test.BodyPostGet,
			},
				grpc.Header(&header))
			require.NoError(t, err)
			assert.True(t, strings.Contains(resp.Addres, test.want.RespBody))
		case "post amd get":
			// тест на запрос POST и получение идентификатора
			// затем получение исходного URL по идентификатору
			fmt.Println(test.name)
			md := metadata.New(map[string]string{"UserID": test.UserID})
			ctx := metadata.NewOutgoingContext(context.Background(), md)

			var header metadata.MD
			// добавляем пользователей
			resp, err := c.PostAddr(ctx, &pb.PostAddrRequset{
				Addres: test.BodyPostGet,
			},
				grpc.Header(&header))
			require.NoError(t, err)
			require.True(t, strings.Contains(resp.Addres, test.want.RespBody))
			shortAddr := resp.Addres
			resp2, err := c.GetAddr(ctx, &pb.GetAddrRequset{
				ShortAddr: shortAddr,
			})
			require.NoError(t, err)
			assert.Equal(t, resp2.Addres, test.BodyPostGet)
		case "batch post":
			// тест на множественный POST для иного пользователя
			// и получение всех сокращенныых URL для этого пользователя
			fmt.Println(test.name)
			md := metadata.New(map[string]string{"UserID": test.UserID})
			ctx := metadata.NewOutgoingContext(context.Background(), md)

			var header metadata.MD
			// добавляем пользователей
			resp, err := c.PostBatch(ctx, &pb.PostBatchRequset{
				BatchURL: test.BodyBatch,
			},
				grpc.Header(&header))
			require.NoError(t, err)
			for i, val := range resp.BatchURL {
				assert.Equal(t, test.want.RespBodyBatch[i].CorrID, val.CorrID)
				assert.True(t, strings.Contains(val.URL, test.want.RespBodyBatch[i].URL))
			}
			resp2, err := c.GetAddrUser(ctx, &pb.GetAddrUserRequset{}, grpc.Header(&header))
			require.NoError(t, err)
			for i, val := range resp2.UserURLs {
				assert.Equal(t, test.BodyBatch[i].URL, val.OrigURL)
				assert.True(t, strings.Contains(val.ShortURL, test.want.RespBodyBatch[i].URL))
			}
		case "bad user":
			// тест на запрос POST с плохим userID
			fmt.Println(test.name)
			md := metadata.New(map[string]string{"UserID": test.UserID})
			ctx := metadata.NewOutgoingContext(context.Background(), md)

			var header metadata.MD
			// добавляем пользователей
			_, err := c.PostAddr(ctx, &pb.PostAddrRequset{
				Addres: test.BodyPostGet,
			},
				grpc.Header(&header))
			require.NoError(t, err)
			NewUserID := header.Get("UserID")
			assert.True(t, len(NewUserID) > 0)
		// case "delete addr":
		//	// ДОСТУПНО ТОЛЬКО С БАЗОЙ
		//	// тест на множественный POST и удаление одного из URL
		//  // затем попытка получить удаленный URL
		// 	fmt.Println(test.name)
		// 	md := metadata.New(map[string]string{"UserID": test.UserID})
		// 	ctx := metadata.NewOutgoingContext(context.Background(), md)

		// 	var header metadata.MD
		// 	// добавляем пользователей
		// 	resp, err := c.PostBatch(ctx, &pb.PostBatchRequset{
		// 		BatchURL: test.BodyBatch,
		// 	},
		// 		grpc.Header(&header))
		// 	require.NoError(t, err)
		// 	var ForDelete []string
		// 	var shortAddrs []string
		// 	for i, val := range resp.BatchURL {
		// 		assert.Equal(t, test.BodyBatch[i].CorrID, val.CorrID)
		// 		assert.True(t, strings.Contains(val.URL, "http://example.com"))
		// 		shortAddrs = append(shortAddrs, val.URL)
		// 		ident, _ := strings.CutPrefix(val.URL, "http://example.com/")
		// 		ForDelete = append(ForDelete, ident)
		// 	}
		// 	fmt.Println(ForDelete)
		// 	_, err = c.DeleteAddr(ctx, &pb.DeleteAddrRequest{
		// 		Idents: ForDelete,
		// 	}, grpc.Header(&header))
		// 	require.NoError(t, err)

		// 	_, err = c.GetAddr(ctx, &pb.GetAddrRequset{
		// 		ShortAddr: shortAddrs[0],
		// 	}, grpc.Header(&header))
		// 	fmt.Println(err)
		// 	assert.True(t, err != nil)
		case "get stats":
			// тест на получение статистики для пользователя
			// из доверенной подсети
			fmt.Println(test.name)
			md := metadata.New(map[string]string{"X-Real-IP": test.UserIP})
			ctx := metadata.NewOutgoingContext(context.Background(), md)

			var header metadata.MD
			// добавляем пользователей
			resp, err := c.GetStats(ctx, &pb.GetStatsRequest{}, grpc.Header(&header))
			require.NoError(t, err)
			assert.Equal(t, test.want.NumUrls, resp.NumURLs)
			assert.Equal(t, test.want.NumUsers, resp.NumUsers)
		}
	}

}
