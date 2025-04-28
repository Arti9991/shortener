package protoServer

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/Arti9991/shortener/internal/app/handlers"
	"github.com/Arti9991/shortener/internal/config"
	"github.com/Arti9991/shortener/internal/gRPC/proto"
	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
	"github.com/Arti9991/shortener/internal/storage/database"
	"github.com/Arti9991/shortener/internal/storage/files"
	"github.com/Arti9991/shortener/internal/storage/inmemory"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
)

type ProtoServer struct {
	Inmemory *inmemory.Data
	Files    *files.FileData
	DataBase *database.DBStor
	Hd       *handlers.HandlersData
	Config   config.Config
	proto.UnimplementedShortenerServer
}

func InitServer() (*ProtoServer, error) {

	var ProtoServ ProtoServer
	var err error
	ctx := context.Background()

	// установка сида для случайных чисел
	rand.Seed(uint64(time.Now().UnixNano()))

	// инициализация конфигурации
	ProtoServ.Config = config.InitConf()
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

// StorInit функция инциализации хранилища с выбором режима хранения (в базе или в памяти).
func (s *ProtoServer) StorInit(ShutDownCtx context.Context, wg *sync.WaitGroup, subIP string) {
	var err1 error
	var err2 error
	// иницализация канала для удаленных URL.
	DeleteOutCh := make(chan models.DeleteURL)
	// инциализация хранилища в базе данных
	s.DataBase, err1 = database.DBinit(s.Config.DBAddress)
	if err1 == nil {
		// ошибка нулевая, работа продолжается через БД.
		// инциализация структуры для файлов.
		s.Files, err2 = files.NewFiles(s.Config.FilePath)
		if err2 != nil {
			logger.Log.Info("Error in creating or file! Setting file or inmemory mode!", zap.Error(err2))
		}
		//инциализируем хранилище данных для хэндлеров с нужным интерфейсом под базу.
		s.Hd = handlers.NewHandlersData(s.DataBase, s.Config.BaseAdr, s.Files, DeleteOutCh, ShutDownCtx, wg, subIP)
		return
	} else {
		//при инцииализации базы возникла ошибка, работа продолжается с внутренней памятью.
		logger.Log.Info("Error while connecting to database! Setting file or inmemory mode!", zap.Error(err1))
		// инциализация структуры для файлов
		s.Files, err2 = files.NewFiles(s.Config.FilePath)
		if err2 != nil {
			logger.Log.Info("Error in creating or file! Setting file or inmemory mode!", zap.Error(err2))
		}
		// инциализация хранилища в памяти.
		s.Inmemory = inmemory.NewData()
		// инциализация хранилища данных для хэндлеров с нужным интерфейсом под память.
		s.Hd = handlers.NewHandlersData(s.Inmemory, s.Config.BaseAdr, s.Files, DeleteOutCh, ShutDownCtx, wg, subIP)
		return
	}
}

func RunGRPCServer() error {
	serv, err := InitServer()
	if err != nil {
		log.Fatal(err)
	}

	// определяем адрес для сервера
	listen, err := net.Listen("tcp", serv.Config.HostAdr)
	if err != nil {
		log.Fatal(err)
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer(grpc.UnaryInterceptor(atuhInterceptor))
	// регистрируем сервис

	proto.RegisterShortenerServer(s, serv)

	logger.Log.Info("New server initialyzed!",
		zap.String("Server addres:", serv.Config.HostAdr),
		zap.String("Base addres:", serv.Config.BaseAdr),
	)
	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
	return nil
}
