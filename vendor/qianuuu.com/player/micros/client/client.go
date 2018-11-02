//
// Author: leafsoar
// Date: 2017-03-16 18:25:38
//

package micros

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"net/url"

	"context"

	"encoding/json"

	"google.golang.org/grpc"
	"qianuuu.com/lib/logs"
	"qianuuu.com/lib/values"
	"qianuuu.com/player/micros/pb"
)

type session struct {
	conn *grpc.ClientConn
	svr  pb.PlayerServiceClient
}

func (s *session) destroy() {
	_ = s.conn.Close()
}

// Client 微服务客户端
type Client struct {
	target string
	opts   []grpc.DialOption
}

// NewClient 创建一个客户端
func NewClient(target string, opts ...grpc.DialOption) *Client {
	ret := &Client{
		target: target,
	}
	for _, item := range opts {
		ret.opts = append(ret.opts, item)
	}
	return ret
}

func (c *Client) newSession() (*session, error) {
	conn, err := grpc.Dial(c.target, c.opts...)
	if err != nil {
		return nil, err
	}

	ret := &session{
		conn: conn,
		svr:  pb.NewPlayerServiceClient(conn),
	}
	return ret, nil
}

// SayHello 测试函数
func (c *Client) SayHello() error {
	ss, err := c.newSession()
	if err != nil {
		return err
	}
	defer ss.destroy()
	r, err := ss.svr.SayHello(context.Background(), &pb.SimpleRequest{StrVal: "leaf"})
	if err != nil {
		return errors.New(url2string(grpc.ErrorDesc(err)))
	}
	logs.Info("rpc hello: %v", r.StrVal)
	return nil
}

// CheckToken token 检测
func (c *Client) CheckToken(uid int, token string) (*pb.SimpleReply, error) {
	ss, err := c.newSession()
	if err != nil {
		return nil, err
	}
	defer ss.destroy()
	request := &pb.SimpleRequest{
		IntVal: int32(uid),
		StrVal: token,
	}
	ret, err := ss.svr.CheckToken(context.Background(), request)
	if err != nil {
		return nil, fmt.Errorf("%s", url2string(grpc.ErrorDesc(err)))
	}
	return ret, nil
}

// UserInfo 用户信息
func (c *Client) UserInfo(uid int) (*pb.User, error) {
	ss, err := c.newSession()
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败 %v", uid)
	}
	defer ss.destroy()
	ret, err := ss.svr.UserInfo(
		context.Background(),
		&pb.SimpleRequest{IntVal: int32(uid)})
	if err != nil {
		return nil, fmt.Errorf("%s", url2string(grpc.ErrorDesc(err)))
	}
	return ret, nil
}

// AddTableRecord 添加录像
func (c *Client) AddTableRecord(appname string, tid, inning int, uids, scores string, data []byte) (*pb.SimpleReply, error) {
	if tid == 0 || inning == 0 {
		return nil, fmt.Errorf("添加 record 记录失败 tid %v 或者 inning %v 不正确！", tid, inning)
	}
	if uids == "" || scores == "" {
		return nil, fmt.Errorf("添加 record 记录失败 tid %v 或者 inning %v 不正确！", tid, inning)
	}
	if len(data) <= 0 {
		return nil, fmt.Errorf("添加 record 记录失败，记录信息不能为空！")
	}
	ss, err := c.newSession()
	if err != nil {
		return nil, fmt.Errorf("添加战绩失败 tid: %v  inning: %v", tid, inning)
	}
	defer ss.destroy()
	zdata, err := func(input []byte) ([]byte, error) {
		var buf bytes.Buffer
		compressor := zlib.NewWriter(&buf)
		if _, err := compressor.Write(data); err != nil {
			return nil, err
		}
		if err := compressor.Close(); err != nil {
			return nil, err

		}
		return buf.Bytes(), nil
	}(data)
	if err != nil {
		return nil, err
	}
	in := &pb.TableRecordRequest{
		Appname: appname,
		Tableid: int32(tid),
		Inning:  int32(inning),
		Uids:    uids,
		Scores:  scores,
		Zdata:   zdata,
	}
	ret, err := ss.svr.AddTableRecord(context.Background(), in)
	if err != nil {
		return nil, fmt.Errorf("%s", url2string(grpc.ErrorDesc(err)))
	}
	return ret, nil
}

// ConsumeGoods 消耗物品 (消耗 id， error)
func (c *Client) ConsumeGoods(uid int, goodtype string, count int, game string, tableid int, parms values.ValueMap) (int, error) {
	if uid == 0 || goodtype == "" || count == 0 || game == "" || tableid == 0 {
		return 0, fmt.Errorf("请传入合法的消耗数据 uid: %v gt: %v count: %v game: %v tableid: %v",
			uid, goodtype, count, game, tableid)
	}
	ss, err := c.newSession()
	if err != nil {
		return 0, fmt.Errorf("消耗操作错误 %v", uid)
	}
	defer ss.destroy()
	mp := parms
	if mp == nil {
		mp = map[string]interface{}{}
	}
	mp["app_name"] = game
	mp["table_id"] = tableid
	detail, _ := json.Marshal(mp)
	consume := &pb.ConsumeRequest{
		Uid:      int32(uid),
		GoodType: goodtype,
		Count:    int32(count),
		Detail:   string(detail),
	}
	rp, err := ss.svr.ConsumeGoods(
		context.Background(), consume)
	if err != nil {
		return 0, fmt.Errorf("%s", url2string(grpc.ErrorDesc(err)))
	}
	return int(rp.IntVal), err
}

// EarnGoods 获得物品
func (c *Client) EarnGoods(uid int, goodtype string, count int, game string, tableid int, parms values.ValueMap) (int, error) {
	if uid == 0 || goodtype == "" || count == 0 || game == "" || tableid == 0 {
		return 0, fmt.Errorf("请传入合法的消耗数据 uid: %v gt: %v count: %v game: %v tableid: %v",
			uid, goodtype, count, game, tableid)
	}
	ss, err := c.newSession()
	if err != nil {
		return 0, fmt.Errorf("奖励操作错误 %v", uid)
	}
	defer ss.destroy()
	//detail := fmt.Sprintf(`{"app_name": "%s", "table_id": %d}`, game, tableid)
	mp := parms
	if mp == nil {
		mp = map[string]interface{}{}
	}
	mp["app_name"] = game
	mp["table_id"] = tableid
	detail, _ := json.Marshal(mp)
	consume := &pb.ConsumeRequest{
		Uid:      int32(uid),
		GoodType: goodtype,
		Count:    int32(count),
		Detail:   string(detail),
	}
	rp, err := ss.svr.EarnGoods(
		context.Background(), consume)
	if err != nil {
		return 0, fmt.Errorf("%s", url2string(grpc.ErrorDesc(err)))
	}
	return int(rp.IntVal), err
}

//消耗服务费
func (c *Client) ChargeCoin(uid int, goodtype string, count int, game string, tableid int, parms values.ValueMap) (int, error) {
	ss, err := c.newSession()
	if err != nil {
		return 0, err
	}
	defer ss.destroy()
	mp := parms
	if mp == nil {
		mp = map[string]interface{}{}
	}
	mp["app_name"] = game
	mp["table_id"] = tableid
	mp["goodtype"] = goodtype
	mp["count"] = count
	detail, _ := json.Marshal(mp)
	Charge := &pb.SimpleRequest{
		IntVal: int32(uid),
		StrVal: string(detail),
	}
	rp, err := ss.svr.ChargeCoin(
		context.Background(), Charge)
	if err != nil {
		return 0, fmt.Errorf("%s", url2string(grpc.ErrorDesc(err)))
	}
	return int(rp.IntVal), nil
}
func (c *Client) WinLossCoin(src int, goodtype string, des int, tid, count int, game string, params values.ValueMap) error {
	if src == 0 || des == 0 || goodtype == "" || count == 0 || game == "" {
		return fmt.Errorf("请传入合法的消耗数据 uid: %v gt: %v count: %v game: %v ",
			src, goodtype, count, game)
	}
	ss, err := c.newSession()
	if err != nil {
		return fmt.Errorf("消耗操作错误 %v", src)
	}
	defer ss.destroy()
	mp := params
	if mp == nil {
		mp = map[string]interface{}{}
	}
	mp["app_name"] = game
	mp["src"] = src
	mp["des"] = des
	mp["table_id"] = tid
	mp["goodtype"] = goodtype
	mp["count"] = count
	detail, _ := json.Marshal(mp)
	WinLoss := &pb.SimpleRequest{
		IntVal: int32(src),
		StrVal: string(detail),
	}
	_, err = ss.svr.WinLossCoin(
		context.Background(), WinLoss)
	if err != nil {
		return fmt.Errorf("%s", url2string(grpc.ErrorDesc(err)))
	}
	return err
}

//MultiConsumeEarnCoin   uids, "coin", counts, svr.appname, tid,parms
func (c *Client) MultiWinLossCoin(uids []int, goodtype string, counts []int, game string, tid int, params values.ValueMap) (int, error) {
	if uids == nil || goodtype == "" || counts == nil || game == "" {
		return 0, fmt.Errorf("请传入合法的消耗数据 uids: %v gt: %v counts: %v game: %v ",
			uids, goodtype, counts, game)
	}
	ss, err := c.newSession()
	if err != nil {
		return 0, fmt.Errorf("消耗操作错误 %v", uids)
	}
	defer ss.destroy()
	mp := params
	if mp == nil {
		mp = map[string]interface{}{}
	}
	mp["app_name"] = game
	mp["table_id"] = tid
	mp["uids"] = uids
	mp["counts"] = counts
	mp["goodtype"] = goodtype
	detail, _ := json.Marshal(mp)
	MultiWinLoss := &pb.SimpleRequest{
		StrVal: string(detail),
	}
	_, err = ss.svr.MultiWinLossCoin(
		context.Background(), MultiWinLoss)
	if err != nil {
		return 0, fmt.Errorf("%s", url2string(grpc.ErrorDesc(err)))
	}
	return int(MultiWinLoss.IntVal), err
}

// ConsumeGoodsLoss 消耗房卡失败 (消耗 id)
func (c *Client) ConsumeGoodsLoss(consumeid int) error {
	if consumeid <= 0 {
		return fmt.Errorf("消耗失败操作错误 %v", consumeid)
	}
	ss, err := c.newSession()
	if err != nil {
		return fmt.Errorf("消耗失败操作错误 %v", consumeid)
	}
	defer ss.destroy()
	_, err = ss.svr.ConsumeGoodsLoss(
		context.Background(),
		&pb.SimpleRequest{IntVal: int32(consumeid)})
	if err != nil {
		return fmt.Errorf("%s", url2string(grpc.ErrorDesc(err)))
	}
	return nil
}

// MultiConsumeGoods 多人消耗房卡
func (c *Client) MultiConsumeGoods(uids []int, goodtype string, count int, game string, tableid int) (int, error) {
	ss, err := c.newSession()
	if err != nil {
		return 0, fmt.Errorf("多人消耗操作错误 %v", uids)
	}

	uids32 := []int32{}
	for _, item := range uids {
		uids32 = append(uids32, int32(item))
	}

	defer ss.destroy()
	detail := fmt.Sprintf(`{"game": "%s", "table_id": %d}`, game, tableid)
	consume := &pb.ConsumeRequest{
		// Uid:      int32(uid),
		GoodType: goodtype,
		Count:    int32(count),
		Detail:   detail,
		Uids:     uids32,
	}
	rp, err := ss.svr.MultiConsumeGoods(
		context.Background(), consume)
	if err != nil {
		return 0, fmt.Errorf("%s", url2string(grpc.ErrorDesc(err)))
	}
	return int(rp.IntVal), err
}

// GetTableID 获取可用桌子 ID
func (c *Client) GetTableID(params values.ValueMap) (int, error) {
	ss, err := c.newSession()
	if err != nil {
		return 0, fmt.Errorf("获取可用桌子 ID 失败")
	}
	defer ss.destroy()
	ret, err := ss.svr.GetTableID(context.Background(), &pb.SimpleRequest{
		StrVal: string(params.ToJSON()),
	})
	if err != nil {
		return 0, fmt.Errorf("%s", url2string(grpc.ErrorDesc(err)))
	}
	if ret != nil {
		return int(ret.IntVal), nil
	}
	return 0, nil
}

// PutTableID 放回可用桌子 ID
func (c *Client) PutTableID(tid int) error {
	// 预留
	return nil
}

func url2string(errmsg string) string {
	ret, _ := url.QueryUnescape(errmsg)
	return ret
}
