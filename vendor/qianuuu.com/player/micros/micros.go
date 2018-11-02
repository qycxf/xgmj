//
// Author: leafsoar
// Date: 2017-03-16 13:58:52
//

package micros

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"qianuuu.com/lib/logs"
	"qianuuu.com/player/cache"
	"qianuuu.com/player/micros/pb"
)

// StartServer 开始一个 player 微服务
func StartServer(laddr string, handler Handler, cache cache.UserCache, jwtSiging string) {
	lis, err := net.Listen("tcp", laddr)
	if err != nil {
		logs.Fatal("failed to listen: %v", err)
		return
	}
	logs.Info("start player microservice " + laddr)
	s := grpc.NewServer()
	pb.RegisterPlayerServiceServer(s, &server{handler: handler, cacheUser: cache, jwtSigning: jwtSiging})
	if err := s.Serve(lis); err != nil {
		fmt.Println(err)
	}
}
