package repository

import (
	"database/sql"
	"github.com/Deansquirrel/goToolMSSql2000"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/object"
)

//门店信道数据
type zb struct {
	db *sql.DB
}

func NewRepZb() (*zb, error) {
	c := common{}
	dbConfig, err := c.getZbDbConfig()
	if err != nil {
		return nil, err
	}
	conn, err := goToolMSSql2000.GetConn(c.convertDbConfigTo2000(dbConfig))
	if err != nil {
		return nil, err
	}
	return &zb{
		db: conn,
	}, nil
}

//获取门店营业日开闭店记录恢复时间（总部库）
func (r *zb) GetMdYyStateTransTime() (*object.MdYyStateTransTimeZb, error) {
	//TODO 获取门店营业日开闭店记录恢复时间（总部库）
	return nil, nil
}
