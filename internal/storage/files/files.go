package files

import (
	"bufio"
	"encoding/json"
	"os"

	"go.uber.org/zap"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
)

// Структура с информацией о файлах.
type FileStor struct {
	Shorturl string `json:"short_url"`
	Origurl  string `json:"original_url"`
	ID       int    `json:"uuid"`
}

// структура для кодирования данных.
type FileData struct {
	Path     string
	ID       int
	InMemory bool
}

// NewFiles конструктор структуры для работы с файлами. Также он создает/проверяет сам файл.
func NewFiles(Path string) (*FileData, error) {
	file, err := os.OpenFile(Path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil || Path == "" {
		return &FileData{InMemory: true}, err
	}
	file.Close()
	return &FileData{ID: 0, Path: Path, InMemory: false}, nil
}

// FilesTest тестовый инциализатор с отключенным флагом.
func FilesTest() *FileData {
	return &FileData{InMemory: true}
}

// FileSave функция сохранения исходного и укороченного URL в файл.
func (d *FileData) FileSave(key string, val string) error {
	// проверка флага на хранение данных в памяти
	if d.InMemory {
		return nil
	}
	logger.Log.Info("INFO Saving file")
	file, err := os.OpenFile(d.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logger.Log.Info("Error in saving file! Setting in memory mode!", zap.Error(err))
		d.InMemory = true
		return err
	}
	writer := bufio.NewWriter(file)
	defer file.Close()
	d.ID += 1
	var fl FileStor
	fl.ID = d.ID
	fl.Origurl = val
	fl.Shorturl = key
	fl.ID = d.ID
	data, err := json.Marshal(&fl)
	if err != nil {
		logger.Log.Info("Error in marshalling data!", zap.Error(err))
	}

	//запись строки
	if _, err := writer.Write(data); err != nil {
		logger.Log.Info("Error in writing data!", zap.Error(err))
	}
	//запись символа переноса строки
	if err := writer.WriteByte('\n'); err != nil {
		logger.Log.Info("Error in writing data!", zap.Error(err))
	}

	if err := writer.Flush(); err != nil {
		logger.Log.Info("Error in flashing data!", zap.Error(err))
	}
	return nil
}

// FileSaveTx функция сохранения множества URL в файл при чтении их из JSON.
func (d *FileData) FileSaveTx(InURLs models.InBuff, BaseAdr string) error {
	// проверка флага на хранение данных в памяти
	if d.InMemory {
		return nil
	}

	logger.Log.Info("INFO Saving file")
	file, err := os.OpenFile(d.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logger.Log.Info("Error in saving file! Setting in memory mode!", zap.Error(err))
		d.InMemory = true
		return err
	}
	writer := bufio.NewWriter(file)
	defer file.Close()

	for _, income := range InURLs {

		hashStr := income.Hash
		d.ID += 1
		var fl FileStor
		fl.ID = d.ID
		fl.Origurl = income.URL
		fl.Shorturl = hashStr

		data, err := json.Marshal(&fl)
		if err != nil {
			logger.Log.Info("Error in marshalling data!", zap.Error(err))
		}

		//запись строки
		if _, err := writer.Write(data); err != nil {
			logger.Log.Info("Error in writing data!", zap.Error(err))
		}
		//запись символа переноса строки
		if err := writer.WriteByte('\n'); err != nil {
			logger.Log.Info("Error in writing data!", zap.Error(err))
		}

		if err := writer.Flush(); err != nil {
			logger.Log.Info("Error in flashing data!", zap.Error(err))
		}

	}
	return nil
}
