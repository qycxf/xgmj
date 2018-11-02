//
// Author: leafsoar
// Date: 2017-03-16 18:20:51
//

package main

import (
	"fmt"

	"qianuuu.com/player"
	mc "qianuuu.com/player/micros/client"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("client ...")
	// url := "54.222.229.175:8387"
	url := "localhost:8507"
	// url := "twxz.qianuuu.cn:8407"
	// url := "test.qianuuu.cn:8507"
	// url = "fymj.yytxgame.com:8357"
	// url = "101.37.16.206:8407"
	client := mc.NewClient(url, grpc.WithInsecure())
	_ = client.SayHello()

	// fmt.Println(client.GetTableID("biji"))
	// fmt.Println(client.UserInfo(10000))
	// fmt.Println(client.EarnGoods(10000, "coin", 3, "twxz", 10000))

	// fmt.Println(client.MultiConsumeGoods([]int{23923, 10000}, "fangka", 1, "fymj", 10000))
	// ret, err := client.EarnGoods(10001, "coin", 10, "nnpk", 10000)
	// fmt.Println(ret, err)
	// for i := 0; i < 10; i++ {
	// fmt.Println(client.GetTableID())
	// }

	// fmt.Println(client.ConsumeGoodsLoss(820))

	// user, err := client.UserInfo(10000)
	// fmt.Printf(" %v %v \n", user, err)
	// ret, err := client.CheckToken(10000, "eyJbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTEzNzY1MTksImlkIjoxMDAwMCwibmFtZSI6IjEwMDAwIn0.SzYoC-woWSjATjsi1U5gRNnnfXYJ7PPUt8mDfaON8xU")
	// fmt.Println(ret, err)

	scores := player.Scores{
		{UID: 101, NickName: "user1", Score: 8},
		{UID: 102, NickName: "user2", Score: 9},
		{UID: 103, NickName: "user3", Score: 10},
		{UID: 104, NickName: "user4", Score: 11},
	}
	data := []byte(`[{"tableCfg":{"tableId":498981,"gameCt":10,"playerCt":2,"present":false,"doubleCt":1,"dai2Jokers":false},"uuid":"uuid","state":0,"seatInfo":[{"seatId":0,"state":1,"playerInfo":{"uid":105451,"nickName":"游客_105451"},"score":[1,0],"offline":false,"canGiveUp":false},{"seatId":1,"state":1,"playerInfo":{"uid":105040,"nickName":"→_→等←_←"},"score":[1,0],"offline":false,"canGiveUp":false}],"reaminTime":0,"gameCt":3,"tableTime":1502422230,"startSeatId":0},{"tableCfg":{"tableId":498981,"gameCt":10,"playerCt":2,"present":false,"doubleCt":1,"dai2Jokers":false},"uuid":"uuid","state":1,"seatInfo":[{"seatId":0,"state":2,"playerInfo":{"uid":105451,"nickName":"游客_105451"},"score":[1,0],"offline":false,"canGiveUp":false},{"seatId":1,"state":1,"playerInfo":{"uid":105040,"nickName":"→_→等←_←"},"score":[1,0],"offline":false,"canGiveUp":false}],"reaminTime":15,"gameCt":3,"tableTime":1502422235,"startSeatId":0},{"tableCfg":{"tableId":498981,"gameCt":10,"playerCt":2,"present":false,"doubleCt":1,"dai2Jokers":false},"uuid":"uuid","state":2,"seatInfo":[{"seatId":0,"state":2,"playerInfo":{"uid":105451,"nickName":"游客_105451"},"handCards":[38,62,30,14,54,34,53,2,12],"score":[1,0],"offline":false,"canGiveUp":true},{"seatId":1,"state":2,"playerInfo":{"uid":105040,"nickName":"→_→等←_←"},"handCards":[42,40,8,6,56,22,9,23,60],"score":[1,0],"offline":false,"canGiveUp":true}],"reaminTime":0,"gameCt":3,"tableTime":1502422236,"startSeatId":0},{"tableCfg":{"tableId":498981,"gameCt":10,"playerCt":2,"present":false,"doubleCt":1,"dai2Jokers":false},"uuid":"uuid","state":3,"seatInfo":[{"seatId":0,"state":2,"playerInfo":{"uid":105451,"nickName":"游客_105451"},"handCards":[38,62,30,14,54,34,53,2,12],"score":[1,0],"offline":false,"canGiveUp":true},{"seatId":1,"state":2,"playerInfo":{"uid":105040,"nickName":"→_→等←_←"},"handCards":[42,40,8,6,56,22,9,23,60],"score":[1,0],"offline":false,"canGiveUp":true}],"reaminTime":10,"gameCt":3,"tableTime":1502422238,"startSeatId":0},{"tableCfg":{"tableId":498981,"gameCt":10,"playerCt":2,"present":false,"doubleCt":1,"dai2Jokers":false},"uuid":"uuid","state":3,"seatInfo":[{"seatId":0,"state":3,"playerInfo":{"uid":105451,"nickName":"游客_105451"},"handCards":[34,2,53,54,38,12,62,30,14],"score":[1,0],"offline":false,"canGiveUp":true},{"seatId":1,"state":2,"playerInfo":{"uid":105040,"nickName":"→_→等←_←"},"handCards":[42,40,8,6,56,22,9,23,60],"score":[1,0],"offline":false,"canGiveUp":true}],"reaminTime":3,"gameCt":3,"tableTime":1502422245,"startSeatId":0},{"tableCfg":{"tableId":498981,"gameCt":10,"playerCt":2,"present":false,"doubleCt":1,"dai2Jokers":false},"uuid":"uuid","state":4,"seatInfo":[{"seatId":0,"state":3,"playerInfo":{"uid":105451,"nickName":"游客_105451"},"handCards":[34,2,53,54,38,12,62,30,14],"score":[1,4],"offline":false,"canGiveUp":true},{"seatId":1,"state":3,"playerInfo":{"uid":105040,"nickName":"→_→等←_←"},"handCards":[60,42,9,22,6,23,56,40,8],"score":[0,4],"offline":false,"canGiveUp":true}],"reaminTime":8,"gameCt":3,"tableTime":1502422249,"seatResult":[{"seatId":0,"pxId":2,"pxName":"对子","isWinner":true,"score":4,"daoScSign":[1,1,1],"daoScore":[1,1,1],"cxInfo":[{"pxId":3,"pxName":"三通","scsign":1,"score":1}]},{"seatId":1,"pxId":2,"pxName":"对子","isWinner":false,"score":-4,"daoScSign":[0,0,0],"daoScore":[1,1,1]}],"startSeatId":0},{"tableCfg":{"tableId":498981,"gameCt":10,"playerCt":2,"present":false,"doubleCt":1,"dai2Jokers":false},"uuid":"uuid","state":5,"seatInfo":[{"seatId":0,"state":3,"playerInfo":{"uid":105451,"nickName":"游客_105451"},"handCards":[34,2,53,54,38,12,62,30,14],"score":[1,4],"offline":false,"canGiveUp":true},{"seatId":1,"state":3,"playerInfo":{"uid":105040,"nickName":"→_→等←_←"},"handCards":[60,42,9,22,6,23,56,40,8],"score":[0,4],"offline":false,"canGiveUp":true}],"reaminTime":0,"gameCt":3,"tableTime":1502422257,"seatResult":[{"seatId":0,"pxId":2,"pxName":"对子","isWinner":true,"score":4,"daoScSign":[1,1,1],"daoScore":[1,1,1],"cxInfo":[{"pxId":3,"pxName":"三通","scsign":1,"score":1}]},{"seatId":1,"pxId":2,"pxName":"对子","isWinner":false,"score":-4,"daoScSign":[0,0,0],"daoScore":[1,1,1]}],"startSeatId":0}]`)
	ret, err := client.AddTableRecord("test", 10000, 4, scores.Uids(), scores.String(), data)
	fmt.Println(ret, err)
}
