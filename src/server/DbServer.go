package server

import (
	"fmt"
	// "strconv"
	// "strings"
    // "regexp"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/jmoiron/sqlx"
	
	// . "dao"
	// . "model"
)

type DbServer struct{
	db *sqlx.DB;
}

var GetDbServer = (func() (func() (*DbServer)) {
	var instance *DbServer;

	return func() (*DbServer) {
		if(instance == nil) {
			instance = new(DbServer);
		}
		return instance;
	};
})();

//数据库
//连接数据库
func (c *DbServer) Init(path string){
	// InitDbDao(path);

	var err error
	c.db, err = sqlx.Connect("sqlite3", path)
	c.checkErr(err)
}

//执行sql
func (c *DbServer) DirectExec(strSql string) bool {
	// _, err := c.db.Exec(strSql);
	
	// return err == nil;

	return c.Exec(strSql, make([]interface{}, 0));
}

//执行sql
func (c *DbServer) Exec(strSql string, arrParam []interface{}) bool {
	_, err := c.db.Exec(strSql, arrParam...);
	
	return err == nil;
}

//插入sql
func (c *DbServer) DirectInsert(strSql string) (int64, error) {
	return c.Insert(strSql, make([]interface{}, 0));
	// rst, err := c.db.Exec(strSql);
	// if(err != nil) {
	// 	return -1, err;
	// }

	// var id int64;
	// id, err = rst.LastInsertId();

	// if(err != nil) {
	// 	return -1, err;
	// }
	
	// return id, nil;
}

//插入sql
func (c *DbServer) Insert(strSql string, arrParam []interface{}) (int64, error) {
	rst, err := c.db.Exec(strSql, arrParam...);
	if(err != nil) {
		return -1, err;
	}

	var id int64;
	id, err = rst.LastInsertId();

	if(err != nil) {
		return -1, err;
	}
	
	return id, nil;
}

//查询sql
func (c *DbServer) DirectQuery(strSql string) []map[string]interface{}{
	return c.Query(strSql, make([]interface{}, 0));
}

//查询sql
func (c *DbServer) Query(strSql string, arrParam []interface{}) []map[string]interface{}{
    rows, err := c.db.Queryx(strSql, arrParam...);
	c.checkErr(err);

	slice := []map[string]interface{}{};
	if(err != nil) {
		return slice;
	}

	var colType []*sql.ColumnType = nil;
	for rows.Next() {
		if(colType == nil) {
			// colName, err = rows.Columns();
			colType, err = rows.ColumnTypes();
		}

		mapData := make(map[string]interface{});

		arr := c.getRowListData(colType);
		err = rows.Scan(arr...);
		c.checkErr(err);
		
		for i := 0; i < len(colType); i++ {
			mapData[colType[i].Name()] = arr[i];
		}
		slice = append(slice, mapData);
	}
	defer rows.Close();

	return slice;
}

func (c *DbServer) getRowListData(colType []*sql.ColumnType) []interface{}{
	arr := make([]interface{}, len(colType));

	for i := 0; i < len(colType); i++ {
		fieldType := colType[i].DatabaseTypeName();
		switch(fieldType) {
		case "INTEGER": {
			var tmp int = 0;
			arr[i] = &tmp;
		}
		case "REAL": {
			var tmp float32 = 0;
			arr[i] = &tmp;
		}
		case "BOLB": {
			var tmp []byte = nil;
			arr[i] = &tmp;
		}
		default: {
			var tmp string = "";
			arr[i] = &tmp;
		}
		}
	}
	return arr;
}

func (c *DbServer) checkErr(err error) {
    if err != nil {
        // panic(err)
        fmt.Println(err)
    }
}
