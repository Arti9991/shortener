package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"net/http"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/models"
	"go.uber.org/zap"
)

var UserSession = "userID"
var key = []byte{183, 21, 219, 229, 199, 223, 64, 207, 94, 48, 138, 6, 9, 250, 124, 17}

func MiddlewareAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var UserExist bool
		var UserID string

		cookie, err := req.Cookie(UserSession)
		if err != nil {
			UserExist = false
		} else {
			UserID, err = DecodeUserID(cookie.Value)
			if err != nil {
				logger.Log.Info("Error in Decoding", zap.Error(err))
				UserExist = false
			} else {
				UserExist = true
			}
		}

		if !UserExist {
			UserID = models.RandomString(16)
			UserEnc, err := EncodeUserID(UserID)
			if err != nil {
				logger.Log.Info("Error in Encoding", zap.Error(err))
			}
			cookie = &http.Cookie{
				Name:  UserSession,
				Value: UserEnc,
			}
			http.SetCookie(res, cookie)
		}
		//fmt.Printf("\n\nUserID in context: %s\n\n", UserID)
		ctx := context.WithValue(req.Context(), models.CtxKey, models.UserInfo{UserID: UserID, Register: UserExist})
		req = req.WithContext(ctx)
		// передаём управление хендлеру
		h.ServeHTTP(res, req)
	}
}

func MakeCiper() (cipher.Block, error) {
	// получаем cipher.Block
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return aesblock, nil
}

func EncodeUserID(UserID string) (string, error) {
	cip, err := MakeCiper()
	if err != nil {
		return "", err
	}
	UserIDenc := make([]byte, aes.BlockSize) // зашифровываем
	cip.Encrypt(UserIDenc, []byte(UserID))
	//fmt.Printf("UserID: %x\n", UserID)
	//fmt.Printf("encrypted: %x\n", UserIDenc)
	return hex.EncodeToString(UserIDenc), nil
}

func DecodeUserID(UserIDenc64 string) (string, error) {
	UserIDenc, err := hex.DecodeString(UserIDenc64)
	if err != nil {
		return "", err
	}
	cip, err := MakeCiper()
	if err != nil {
		return "", err
	}
	UserID := make([]byte, aes.BlockSize) // расшифровываем
	cip.Decrypt(UserID, UserIDenc)
	//fmt.Printf("decrypted: %x\n", UserID)
	return string(UserID), nil
}
