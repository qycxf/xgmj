//
// Author: leafsoar
// Date: 2017-03-16 15:46:34
//

package main

import (
	"fmt"

	_ "github.com/lib/pq"

	"qianuuu.com/lib/client"
	"qianuuu.com/player/cache"
	"qianuuu.com/player/login/handler"
	"qianuuu.com/player/micros"
	"qianuuu.com/player/usecase"
)

// UCDBURL URL
const UCDBURL = "postgres://postgres:qianuuu_12345@test.qianuuu.cn:5462/channel?sslmode=disable"

func main() {
	client, err := client.NewClient("postgres", UCDBURL)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	uc := usecase.NewUsecase(client)
	defer func() {
		_ = client.Close()
	}()
	_ = uc

	handler := &handler.HandlerImpl{
		Uc: uc,
	}

	micros.StartServer(":8080", handler, cache.NewUserCache("", 1), "jwt_signing")
}
