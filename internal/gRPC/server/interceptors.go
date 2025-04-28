package protoServer

import (
	"context"
	"fmt"

	"github.com/Arti9991/shortener/internal/app/auth"
	"github.com/Arti9991/shortener/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func atuhInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	var UserID string
	var err error
	UserExist := true

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("UserID")
		if len(values) > 0 {
			UserIDhash := values[0]
			fmt.Println(len(UserIDhash))
			if len(UserIDhash) != 32 {
				UserExist = false
				newCtx := context.WithValue(ctx, models.CtxKey, models.UserInfo{Register: UserExist})
				return handler(newCtx, req)
			}
			UserID, err = auth.DecodeUserID(UserIDhash)
			if err != nil {
				UserExist = false
			}
			fmt.Println("User ID un interceptor is:", UserID)
		} else if len(values) == 0 {
			UserExist = false
		}
	} else {
		UserExist = false
	}

	newCtx := context.WithValue(ctx, models.CtxKey, models.UserInfo{UserID: UserID, Register: UserExist})
	return handler(newCtx, req)
}
