package table

import (
	"fmt"
	"runtime"

	"qianuuu.com/ahmj/internal/config"
	"qianuuu.com/lib/logs"
	"qianuuu.com/lib/util"
	"qianuuu.com/player"
)

// Tables 桌子管理
type Tables struct {
	tables     *util.Map
	Robots     *Robots
	MsgHandler Handler
}

//初始化 Tables
func InitTables(_msgHandler Handler) *Tables {
	ret := &Tables{
		tables:     &util.Map{},
		Robots:     NewRobots(),
		MsgHandler: _msgHandler,
	}
	return ret
}

// ReadRange 桌子遍历
func (ts *Tables) ReadRange(fn func(tid int, table *Table)) {
	tbs := make([]*Table, 0, ts.tables.Len())
	ts.tables.RLockRange(func(k interface{}, t interface{}) {
		tid := k.(int)
		switch t.(type) {
		case *HFTable:
			tbs = append(tbs, t.(*HFTable).Table)
		case *HZTable:
			tbs = append(tbs, t.(*HZTable).Table)
		case *FYTable:
			tbs = append(tbs, t.(*FYTable).Table)
		case *BBTable:
			tbs = append(tbs, t.(*BBTable).Table)
		case *HYTable:
			tbs = append(tbs, t.(*HYTable).Table)
		default:
			logs.Info(" ---------------> ReadRange GetTable table not found _tableID:%v", tid)
		}
	})
	for _, t := range tbs {
		if t != nil {
			fn(t.ID, t)
		}
	}
}

// GetTable 根据 _tableID 返回桌子
func (ts *Tables) GetTable(_tableID int) *Table {
	t := ts.tables.Get(_tableID)
	switch t.(type) {
	case *HFTable:
		return t.(*HFTable).Table
	case *HZTable:
		return t.(*HZTable).Table
	case *FYTable:
		return t.(*FYTable).Table
	case *BBTable:
		return t.(*BBTable).Table
	case *HYTable:
		return t.(*HYTable).Table
	default:
		logs.Info("  ---------------> GetTable table not found _tableID:%v", _tableID)
	}
	return nil
}

// 创建合肥麻将桌子
func (ts *Tables) CreateHFTable(_tableCfg *config.TableCfg) *HFTable {

	tableId := _tableCfg.TableId
	logs.Info("tableId: %v--------------->创建新的 [合肥麻将] 桌子  ", tableId)

	// 创建并启动一个桌子
	table := NewHFTable(tableId, ts.Robots, _tableCfg)
	table.handler = ts.MsgHandler

	go func() {
		defer func() {
			if err := recover(); err != nil {
				trace := make([]byte, 1<<16)
				n := runtime.Stack(trace, true)
				logs.Info("tableId:%v----------->%v", tableId, fmt.Errorf("panic recover\n %v\n stack trace %d bytes\n %s",
					err, n, trace[:n]))
				ts.tables.Del(tableId)
			}
		}()
		table.Serve()
	}()

	//添加到tables中
	ts.tables.Set(tableId, table)
	return table
}

// 创建红中麻将桌子
func (ts *Tables) CreateHZTable(_tableCfg *config.TableCfg) *HZTable {

	tableId := _tableCfg.TableId
	_tableCfg.TableId = tableId
	logs.Info("tableId: %v--------------->创建新的 [红中麻将] 桌子  ", tableId)

	// 创建并启动一个桌子
	table := NewHZTable(tableId, ts.Robots, _tableCfg)
	table.handler = ts.MsgHandler

	go func() {
		defer func() {
			if err := recover(); err != nil {
				trace := make([]byte, 1<<16)
				n := runtime.Stack(trace, true)
				logs.Info("tableId:%v----------->%v", tableId, fmt.Errorf("panic recover\n %v\n stack trace %d bytes\n %s",
					err, n, trace[:n]))
				ts.tables.Del(tableId)
			}
		}()
		table.Serve()
	}()

	//添加到tables中
	ts.tables.Set(tableId, table)

	//logs.Info("============================> len(ts.tables):%v ", ts.tables.Len())
	return table
}

// 创建阜阳麻将桌子
func (ts *Tables) CreateFYTable(_tableCfg *config.TableCfg) *FYTable {

	//tableId := ts.CreateTableID()
	tableId := _tableCfg.TableId
	logs.Info("tableId: %v--------------->创建新的 [阜阳麻将] 桌子  ", tableId)

	// 创建并启动一个桌子
	table := NewFYTable(tableId, ts.Robots, _tableCfg)
	table.handler = ts.MsgHandler

	go func() {
		defer func() {
			if err := recover(); err != nil {
				trace := make([]byte, 1<<16)
				n := runtime.Stack(trace, true)
				logs.Info("tableId:%v----------->%v", tableId, fmt.Errorf("panic recover\n %v\n stack trace %d bytes\n %s",
					err, n, trace[:n]))
				ts.tables.Del(tableId)
			}
		}()
		table.Serve()
	}()

	//添加到tables中
	ts.tables.Set(tableId, table)

	//logs.Info("============================> len(ts.tables):%v ", ts.tables.Len())
	return table
}

// 创建蚌埠麻将桌子
func (ts *Tables) CreateBBTable(_tableCfg *config.TableCfg) *BBTable {

	tableId := _tableCfg.TableId
	//tableId := ts.CreateTableID()
	_tableCfg.TableId = tableId
	logs.Info("tableId: %v--------------->创建新的 [蚌埠麻将] 桌子  ", tableId)

	// 创建并启动一个桌子
	table := NewBBTable(tableId, ts.Robots, _tableCfg)
	table.handler = ts.MsgHandler

	go func() {
		defer func() {
			if err := recover(); err != nil {
				trace := make([]byte, 1<<16)
				n := runtime.Stack(trace, true)
				logs.Info("tableId:%v----------->%v", tableId, fmt.Errorf("panic recover\n %v\n stack trace %d bytes\n %s",
					err, n, trace[:n]))
				ts.tables.Del(tableId)
			}
		}()
		table.Serve()
	}()

	//添加到tables中
	ts.tables.Set(tableId, table)

	//logs.Info("============================> len(ts.tables):%v ", ts.tables.Len())
	return table
}

// 创建怀远麻将桌子
func (ts *Tables) CreateHYTable(_tableCfg *config.TableCfg) *HYTable {

	tableId := _tableCfg.TableId
	//_tableCfg.TableId = tableId
	logs.Info("tableId: %v--------------->创建新的 [怀远麻将] 桌子  ", tableId)

	// 创建并启动一个桌子
	table := NewHYTable(tableId, ts.Robots, _tableCfg)
	table.handler = ts.MsgHandler

	go func() {
		defer func() {
			if err := recover(); err != nil {
				trace := make([]byte, 1<<16)
				n := runtime.Stack(trace, true)
				logs.Info("tableId:%v----------->%v", tableId, fmt.Errorf("panic recover\n %v\n stack trace %d bytes\n %s",
					err, n, trace[:n]))
				ts.tables.Del(tableId)
			}
		}()
		table.Serve()
	}()

	//添加到tables中
	ts.tables.Set(tableId, table)

	//logs.Info("============================> len(ts.tables):%v ", ts.tables.Len())
	return table
}

// 检测是否删除桌子
func (ts *Tables) RemoveTable(_tableId int) {

	table := ts.GetTable(_tableId)
	len1 := ts.tables.Len()
	table.Destroy()
	ts.tables.Del(_tableId)
	len2 := ts.tables.Len()
	logs.Info("tableId:%v---------------RemoveTable------------------------------>len1:%v,len2:%v", _tableId, len1, len2)

}

// 玩家离开牌桌
func (ts *Tables) ExitTable(_player *player.Player) bool {
	if _player.GetTableID() > 0 {
		_table := ts.GetTable(_player.GetTableID())
		if _table != nil {

			canExit := _table.ChkExitTable(_player)
			if canExit {
				_table.Exit(_player) //从桌子上离开
				return true
			} else {
				return false //不能退出
			}

		} else {
			logs.Info("******************** game error !  _table is nil !  _player.TableId:%v", _player.GetTableID())
			return true
		}
	}
	return true
}

//获取桌子总数量
func (ts *Tables) TableCount() int {
	return ts.tables.Len()
}
