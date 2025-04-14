package server

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"

	"github.com/Arti9991/shortener/internal/app/handlers"
	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
	"github.com/Arti9991/shortener/internal/storage/database"
	"github.com/Arti9991/shortener/internal/storage/files"
	"github.com/Arti9991/shortener/internal/storage/inmemory"
)

// StorInit функция инциализации хранилища с выбором режима хранения (в базе или в памяти).
func (s *Server) StorInit(ShutDownCtx context.Context, wg *sync.WaitGroup) {
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
		s.hd = handlers.NewHandlersData(s.DataBase, s.Config.BaseAdr, s.Files, DeleteOutCh, ShutDownCtx, wg)
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
		s.hd = handlers.NewHandlersData(s.Inmemory, s.Config.BaseAdr, s.Files, DeleteOutCh, ShutDownCtx, wg)
		return
	}
}

// FileRead функция для чтения всех данных в файле и сохранения их в базу или память.
func (s *Server) FileRead(d *files.FileData) error {
	// проверка флага на хранение данных в памяти.
	if d.InMemory {
		return nil
	}
	var id int
	logger.Log.Info("INFO reading file")
	file, err := os.OpenFile(d.Path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Info("Error in reading file! Setting in memory mode!", zap.Error(err))
		d.InMemory = true
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		var fl files.FileStor
		buff, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil && err != io.EOF {
			logger.Log.Info("Error in reading data!", zap.Error(err))
			return err
		}
		err = json.Unmarshal(buff, &fl)
		if err != nil {
			logger.Log.Info("Error in unmarshalling data!", zap.Error(err))
			return err
		}

		err = s.hd.Dt.Save(fl.Shorturl, fl.Origurl, "1")
		if err != nil {
			if !strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
				logger.Log.Info("Error in saving data!", zap.Error(err))
				return err
			}
			logger.Log.Info("This URL already in DB!")
		} else {
			id = fl.ID
		}
	}
	d.ID = id
	return nil
}

// RunDeleteStor функция с горутиной, получающей URL для удаления
// из канала и отправки запроса в БД.
func RunDeleteStor(hd *handlers.HandlersData) {
	go func() {
		for {
			select {
			case DelStruct := <-hd.OutDelCh:
				err := hd.Dt.Delete(DelStruct.ShortURL, DelStruct.UserID)
				if err != nil {
					logger.Log.Info("Error in deleting in DB", zap.Error(err))
					continue
				}
			case <-hd.Ctx.Done():
				return
			default:
				continue
			}
		}
	}()
}

// RunWaitShutDown функция для ожидания сигнала о завершении работы сервера
func RunWaitShutDown(hd *handlers.HandlersData, server *http.Server) {
	go func() {
		<-hd.Ctx.Done()
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		logger.Log.Info("Graceful shutdown...")
		if err := server.Shutdown(context.Background()); err != nil {
			// ошибки закрытия Listener
			logger.Log.Info("Error in HTTP server Shutdown: %v", zap.Error(err))
		}
		hd.Wg.Wait()
	}()
}
