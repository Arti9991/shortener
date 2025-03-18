package storage

import (
	"os"
	"testing"

	"github.com/Arti9991/shortener/internal/models"
	"github.com/Arti9991/shortener/internal/storage/files"
	"github.com/Arti9991/shortener/internal/storage/inmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkStor(b *testing.B) {
	//DBInterface, err := database.DBinit("host=localhost user=myuser password=123456 dbname=ShortURL sslmode=disable")
	//require.NoError(b, err)

	FileInterface, err := files.NewFiles("./test_file.csv")
	require.NoError(b, err)

	MemInterface := inmemory.NewData()

	numReps := 400
	hashStr := make([]string, 0)
	URLStr := make([]string, 0)
	UserID := models.RandomString(16)

	for range numReps {
		hashStr = append(hashStr, models.RandomString(8))
		URLStr = append(URLStr, models.RandomString(16))
	}

	//fmt.Printf("\n%s", hashStr)
	// b.Run("STORAGE in database save", func(b *testing.B) {
	// 	for i := 0; i < numReps; i++ {
	// 		DBInterface.Save(hashStr[i], URLStr[i], UserID)
	// 	}
	// })

	// b.Run("STORAGE in database get", func(b *testing.B) {
	// 	for i := 0; i < numReps; i++ {
	// 		res, err := DBInterface.Get(hashStr[i])
	// 		require.NoError(b, err)
	// 		assert.Equal(b, res, URLStr[i])
	// 	}
	// })
	// err = DBInterface.DropTable()
	// require.NoError(b, err)

	b.Run("STORAGE in memory save", func(b *testing.B) {
		for i := 0; i < numReps; i++ {
			err := MemInterface.Save(hashStr[i], URLStr[i], UserID)
			require.NoError(b, err)
		}
	})

	b.Run("STORAGE in memory get", func(b *testing.B) {
		for i := 0; i < numReps; i++ {
			res, err := MemInterface.Get(hashStr[i])
			require.NoError(b, err)
			assert.Equal(b, res, URLStr[i])
		}
	})

	b.Run("STORAGE in file save", func(b *testing.B) {
		for i := 0; i < numReps; i++ {
			err := FileInterface.FileSave(hashStr[i], URLStr[i])
			require.NoError(b, err)
		}
	})

	err = os.Remove("./test_file.csv")
	require.NoError(b, err)

}
