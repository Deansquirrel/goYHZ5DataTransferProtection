package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Deansquirrel/goToolMSSql"
	"github.com/Deansquirrel/goToolMSSql2000"
	"github.com/Deansquirrel/goYHZ5DataTransferProtection/global"
	"strconv"
	"strings"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

type common struct {
}

func NewCommon() *common {
	return &common{}
}

//获取配置库连接配置
func (c *common) GetConfigDbConfig() *goToolMSSql.MSSqlConfig {
	return &goToolMSSql.MSSqlConfig{
		Server: global.SysConfig.DB.Server,
		Port:   global.SysConfig.DB.Port,
		DbName: global.SysConfig.DB.DbName,
		User:   global.SysConfig.DB.User,
		Pwd:    global.SysConfig.DB.Pwd,
	}
}

//将普通数据库连接配置转换为Sql2000可用的配置
func (c *common) ConvertDbConfigTo2000(dbConfig *goToolMSSql.MSSqlConfig) *goToolMSSql2000.MSSqlConfig {
	return &goToolMSSql2000.MSSqlConfig{
		Server: dbConfig.Server,
		Port:   dbConfig.Port,
		DbName: dbConfig.DbName,
		User:   dbConfig.User,
		Pwd:    dbConfig.Pwd,
	}
}

//根据字符串配置，获取数据库连接配置
func (c *common) getDBConfigByStr(connStr string) (*goToolMSSql.MSSqlConfig, error) {
	log.Debug(fmt.Sprintf("根据字符串配置，获取数据库连接配置,str: %s", connStr))
	connStr = strings.Trim(connStr, " ")
	strList := strings.Split(connStr, "|")
	if len(strList) != 5 {
		err := errors.New(fmt.Sprintf("db config num error,exp 5,act %d", len(strList)))
		log.Error(err.Error())
		return nil, err
	}

	port, err := strconv.Atoi(strList[1])
	if err != nil {
		errLs := errors.New(fmt.Sprintf("db config port[%s] trans err: %s", strList[1], err.Error()))
		log.Error(errLs.Error())
		return nil, errLs
	}

	return &goToolMSSql.MSSqlConfig{
		Server: strList[0],
		Port:   port,
		User:   strList[2],
		Pwd:    strList[3],
		DbName: strList[4],
	}, nil
}

func (c *common) GetRowsBySQL(dbConfig *goToolMSSql.MSSqlConfig, sql string, args ...interface{}) (*sql.Rows, error) {
	conn, err := goToolMSSql.GetConn(dbConfig)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if args == nil {
		rows, err := conn.Query(sql)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		return rows, nil
	} else {
		rows, err := conn.Query(sql, args...)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		return rows, nil
	}
}

func (c *common) SetRowsBySQL(dbConfig *goToolMSSql.MSSqlConfig, sql string, args ...interface{}) error {
	conn, err := goToolMSSql.GetConn(dbConfig)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if args == nil {
		_, err = conn.Exec(sql)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		return nil
	} else {
		_, err := conn.Exec(sql, args...)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		return nil
	}
}

func (c *common) GetRowsBySQL2000(dbConfig *goToolMSSql2000.MSSqlConfig, sql string, args ...interface{}) (*sql.Rows, error) {
	conn, err := goToolMSSql2000.GetConn(dbConfig)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if args == nil {
		rows, err := conn.Query(sql)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		return rows, nil
	} else {
		rows, err := conn.Query(sql, args...)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		return rows, nil
	}
}

func (c *common) SetRowsBySQL2000(dbConfig *goToolMSSql2000.MSSqlConfig, sql string, args ...interface{}) error {
	conn, err := goToolMSSql2000.GetConn(dbConfig)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if args == nil {
		_, err = conn.Exec(sql)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		return nil
	} else {
		_, err := conn.Exec(sql, args...)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		return nil
	}
}

//返回默认时间
func (c *common) GetDefaultOprTime() time.Time {
	return time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local)
}
