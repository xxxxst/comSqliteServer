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

type DbServer struct {
	db *sqlx.DB;
}

var GetDbServer = (func() (func(proj string, path string) (*DbServer)) {
	// var instance *DbServer;
	mapServer := make(map[string] *DbServer);

	return func(proj string, path string) (*DbServer) {
		_, ok := mapServer[proj];
		if !ok {
			md := new(DbServer);
			md.Init(path);
			mapServer[proj] = md;
		}
		
		// if(instance == nil) {
		// 	instance = new(DbServer);
		// }
		return mapServer[proj];
	};
})();

//database
//connect database
func (c *DbServer) Init(path string){
	// InitDbDao(path);

	var err error
	c.db, err = sqlx.Connect("sqlite3", path)
	c.checkErr(err)
}

//exec sql
func (c *DbServer) DirectExec(strSql string) bool {
	return c.Exec(strSql, make([]interface{}, 0));
}

//exec sql
func (c *DbServer) Exec(strSql string, arrParam []interface{}) bool {
	_, err := c.db.Exec(strSql, arrParam...);
	
	return err == nil;
}

//insert sql
func (c *DbServer) DirectInsert(strSql string) (int64, error) {
	return c.Insert(strSql, make([]interface{}, 0));
}

//insert sql
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

//query sql
func (c *DbServer) DirectQuery(strSql string) []map[string]interface{}{
	return c.Query(strSql, make([]interface{}, 0));
}

//query sql
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
