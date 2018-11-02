package protobuf

type TableInfoRec struct {
	TableInfoArr []TableInfo
	TableId      int
	Sequence     int
}

func NewTableInfoRec(_tableId int, _sequence int) *TableInfoRec {

	ret := &TableInfoRec{
		TableInfoArr: make([]TableInfo, 0),
		TableId:      _tableId,
		Sequence:     _sequence,
	}
	return ret
}

//添加一条记录
func (tr *TableInfoRec) AddInfoRec(_tableInfo TableInfo) {
	tr.TableInfoArr = append(tr.TableInfoArr, _tableInfo)
}
