//
// Author: leafsoar
// Date: 2017-03-16 15:50:40
//

package micros

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"

	"golang.org/x/net/context"
	// "context"

	"qianuuu.com/lib/logs"
	"qianuuu.com/lib/values"
	"qianuuu.com/player/cache"
	"qianuuu.com/player/login/plugin/center"
	"qianuuu.com/player/micros/pb"
)

type server struct {
	handler    Handler
	jwtSigning string
	cacheUser  cache.UserCache
}

func (s *server) SayHello(ctx context.Context, in *pb.SimpleRequest) (*pb.SimpleReply, error) {
	logs.Info("player microservice sayhello")
	return &pb.SimpleReply{StrVal: "hello " + in.StrVal}, nil
}

func (s *server) GenerateTableID(ctx context.Context, in *pb.SimpleRequest) (*pb.SimpleReply, error) {
	logs.Info("player microservice generate tableid")
	id, err := getTableID()
	return &pb.SimpleReply{StrVal: "hello " + in.StrVal, IntVal: int32(id)}, err
}

func (s *server) CheckToken(ctx context.Context, in *pb.SimpleRequest) (*pb.SimpleReply, error) {
	logs.Info("player microservice checktoken %v", in.IntVal)
	t, err := jwt.Parse(in.StrVal, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSigning), nil
	})
	if err != nil {
		return &pb.SimpleReply{StrVal: err.Error()}, err
	}
	if t == nil || !t.Valid {
		return &pb.SimpleReply{}, errors.New("token 验证失败")
	}
	var vm values.ValueMap
	vm = t.Claims
	if vm.GetInt("id") != int(in.IntVal) {
		return &pb.SimpleReply{}, errors.New("token 验证失败 id 不正确")
	}
	return &pb.SimpleReply{StrVal: "succeed"}, nil
}

func (s *server) UserInfo(ctx context.Context, in *pb.SimpleRequest) (*pb.User, error) {
	logs.Info("player microservice userinfo %v", in.IntVal)
	user, err := s.cacheUser.Get(int(in.IntVal))
	if user == nil || err != nil {
		// user, err = s.uc.UserInfo(int(in.IntVal))
		user, err = s.handler.UserInfo(int(in.IntVal))
		_ = s.cacheUser.Set(int(in.IntVal), user)
	}
	if err != nil {
		return &pb.User{}, err
	}
	return &pb.User{Uid: int32(user.ID), Nickname: user.NickName, Goods: user.Goods}, nil
}

func (s *server) AddTableRecord(ctx context.Context, in *pb.TableRecordRequest) (*pb.SimpleReply, error) {
	info := fmt.Sprintf("tid: %d innging: %d uids: %s scores: %s", in.Tableid, in.Inning, in.Uids, in.Scores)
	logs.Info("player microservice add_tablerecord app: %v tid: %v  inning: %v datalen: %v data:%s",
		in.Appname, in.Tableid, in.Inning, len(in.Zdata), info)
	go func() {
		if !s.handler.EnableOSS() {
			// 如果不使用存储，直接保存战绩到数据库
			_, err := s.handler.AddRecord(int(in.Tableid), int(in.Inning), in.Uids, in.Scores, in.Zdata, in.Appname, "")
			if err != nil {
				logs.Error("add tablerecord err : %s %v", err.Error(), info)
			}
		} else {
			// 将战绩保存到云存储
			d, err := s.unzip(in.Zdata)
			// logs.Info(string(d), err)
			path := ""
			if err == nil {
				path1, err := s.handler.PutTableRecord(d)
				if err != nil {
					logs.Error("table record oss:", err.Error())
				}
				path = path1
			}
			_, err = s.handler.AddRecord(int(in.Tableid), int(in.Inning), in.Uids, in.Scores, nil, in.Appname, path)
			if err != nil {
				logs.Error("add tablerecord err : %s %v", err.Error(), info)
			}
		}
	}()
	return &pb.SimpleReply{IntVal: int32(0)}, nil
}

func (s *server) unzip(content []byte) ([]byte, error) {
	b := bytes.NewReader(content)
	var out bytes.Buffer
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(&out, r)
	return out.Bytes(), err
}

func (s *server) ConsumeGoods(ctx context.Context, in *pb.ConsumeRequest) (*pb.SimpleReply, error) {
	logs.Info("player microservice consume_goods detail: %v", in)
	// cid, err := s.uc.AutoConsumeFangka(int(in.Uid), int(in.Count), in.Detail)
	cid := 0
	var err error
	if in.GoodType == "coin" {
		cid, err = s.handler.ConsumeCoin(int(in.Uid), in.GoodType, int(in.Count), in.Detail)
	} else {
		cid, err = s.handler.AutoConsumeFangka(int(in.Uid), int(in.Count), in.Detail)
	}

	_ = s.cacheUser.Del(int(in.Uid))
	if err != nil {
		logs.Error(err.Error())
		return &pb.SimpleReply{}, err
	}
	return &pb.SimpleReply{IntVal: int32(cid)}, err
}
func (s *server) ChargeCoin(ctx context.Context, in *pb.SimpleRequest) (*pb.SimpleReply, error) {
	logs.Info("player microservice consume_goods detail: %v", in)
	cid, err := s.handler.ChargeCoin(int(in.IntVal), in.StrVal)
	_ = s.cacheUser.Del(int(in.IntVal))
	if err != nil {
		logs.Error(err.Error())
		return &pb.SimpleReply{}, err
	}
	return &pb.SimpleReply{IntVal: int32(cid)}, err
}
func (s *server) WinLossCoin(ctx context.Context, in *pb.SimpleRequest) (*pb.SimpleReply, error) {
	logs.Info("player microservice consume_goods detail: %v", in)
	err := s.handler.WinLossCoin(int(in.IntVal), in.StrVal)
	_ = s.cacheUser.Del(int(in.IntVal))
	if err != nil {
		logs.Error(err.Error())
		return &pb.SimpleReply{}, err
	}
	return &pb.SimpleReply{IntVal: int32(0)}, err
}
func (s *server) MultiWinLossCoin(ctx context.Context, in *pb.SimpleRequest) (*pb.SimpleReply, error) {
	logs.Info("player microservice consume_goods detail: %v", in)
	cid, err := s.handler.MultiWinLossCoin(in.StrVal)
	_ = s.cacheUser.Del(int(in.IntVal))
	if err != nil {
		logs.Error(err.Error())
		return &pb.SimpleReply{}, err
	}
	return &pb.SimpleReply{IntVal: int32(cid)}, err
}
func (s *server) ConsumeGoodsLoss(ctx context.Context, in *pb.SimpleRequest) (*pb.SimpleReply, error) {
	logs.Info("player microservice consume_goods_loss %v", in.IntVal)
	uid := int(in.IntVal)
	// err := s.uc.ConsumeGoodsLoss(uid)
	err := s.handler.ConsumeGoodsLoss(uid)
	_ = s.cacheUser.Del(uid)
	if err != nil {
		logs.Error(err.Error())
	}
	return &pb.SimpleReply{StrVal: "succeed"}, err
}

func (s *server) MultiConsumeGoods(ctx context.Context, in *pb.ConsumeRequest) (*pb.SimpleReply, error) {
	logs.Info("player microservice multi_consume_goods uids:%v detail: %v", in.Uids, in.Detail)
	retid := 0
	var reterr error
	for _, uid32 := range in.Uids {
		// cid, err := s.uc.AutoConsumeFangka(int(uid32), int(in.Count), in.Detail)
		cid, err := s.handler.AutoConsumeFangka(int(uid32), int(in.Count), in.Detail)
		uid := int(uid32)
		_ = s.cacheUser.Del(uid)
		if err != nil {
			logs.Error(err.Error())
		}
		if retid == 0 {
			retid = cid
		}
		if err != nil {
			reterr = err
		}
	}
	return &pb.SimpleReply{IntVal: int32(retid)}, reterr
}

func (s *server) EarnGoods(ctx context.Context, in *pb.ConsumeRequest) (*pb.SimpleReply, error) {
	logs.Info("[microservice] earn_goods detail: %v", in)
	// _, err := s.uc.EarnGoods(int(in.Uid), in.GoodType, int(in.Count), in.Detail)
	_, err := s.handler.EarnGoods(int(in.Uid), in.GoodType, int(in.Count), in.Detail)
	_ = s.cacheUser.Del(int(in.Uid))
	if err != nil {
		logs.Error(err.Error())
	}
	return &pb.SimpleReply{IntVal: 0}, err
}

func (s *server) GetTableID(ctx context.Context, in *pb.SimpleRequest) (*pb.SimpleReply, error) {
	if in.StrVal == "bad" {
		return &pb.SimpleReply{IntVal: int32(0)}, fmt.Errorf("牌桌创建失败，正在维护中，请稍后再试")
	}
	tid, _ := getTableID()
	logs.Info("[microservice] get app: %v tableid : %v", in.StrVal, tid)
	params, err := values.NewValuesFromJSON([]byte(in.GetStrVal()))
	if err != nil {
		// 如果不是 json 数据，直接字符串传入 appanme
		err = center.SetTable(tid, in.StrVal)
	} else {
		err = center.SetTableParams(tid, params)
	}
	if err != nil {
		logs.Error("set table info err: %v", err)
	}
	if tid > 0 {
		return &pb.SimpleReply{IntVal: int32(tid)}, nil
	}
	return &pb.SimpleReply{IntVal: int32(0)}, fmt.Errorf("获取 tableid 失败")
}

func (s *server) PutTableID(ctx context.Context, in *pb.SimpleRequest) (*pb.SimpleReply, error) {
	logs.Info("[microservice] put app: %v tableid : %v", in.StrVal, in.IntVal)
	return &pb.SimpleReply{IntVal: 0}, nil
}
