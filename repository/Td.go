package repository

import (
	"database/sql"
	"github.com/Deansquirrel/goToolMSSql"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/object"
)

//门店信道数据
type td struct {
	db *sql.DB
}

func NewRepTd() (*td, error) {
	c := common{}
	dbConfig, err := c.getTdDbConfig()
	if err != nil {
		return nil, err
	}
	conn, err := goToolMSSql.GetConn(dbConfig)
	if err != nil {
		return nil, err
	}
	return &td{
		db: conn,
	}, nil
}

//获取门店数据传递信道状态
func (r *td) GetMdDataTransState() (*object.MdDataTransState, error) {
	//TODO 获取门店数据传递信道状态
	return nil, nil
}

//获取门店营业日开闭店记录传递时间（通道库）
func (r *td) GetMdYyStateTransTime() (*object.MdYyStateTransTimeTd, error) {
	//TODO 获取门店营业日开闭店记录传递时间（通道库）
	return nil, nil
}
