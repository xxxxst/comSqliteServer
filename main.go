package main

import (
	// "util"
	. "comSqliteServer"
)

func main() {
	// aaa := util.Aaa{};
	// aaa.Show();
	
	srv := GetServer();

	srv.Run();
}
