package comSqliteServer

import (
	"fmt"
	"flag"
    "net"
    "net/http"
	"strconv"
	"time"
	// "log"
	"os"
	// "io/ioutil"
	// "strings"
	// "database/sql"
    "path/filepath"

	// . "dao"
	. "server"
	. "router"
	. "control"
	. "model"
	// util "util"
)

var instance *ComSqliteServer;

type ComSqliteServer struct{
	comMd	*ComModel;
	router	MainRouter;

	cfgCtl	ConfigCtl;
}

func GetServer() (*ComSqliteServer){
	if(instance == nil){
		instance = new(ComSqliteServer);
	}
	return instance;
}

func getStringNotEmpty(val *string, defVal string) string {
	if(val == nil || *val == ""){
		//default path
		return defVal;
	} else{
		return *val;
	}
}

func (c *ComSqliteServer) Run() {
	//show version
	c.comMd = GetComModel();
	c.comMd.Version = "v1.0.0";
	fmt.Println("ComSqliteServer " + c.comMd.Version);

	//get arguments
	configPath := flag.String("configPath", "", "config path");
	flag.Parse();

	//set path
	exeDir, _ := filepath.Abs(filepath.Dir(os.Args[0]));
	c.comMd.ExePath = exeDir + "/";
	c.comMd.ConfigPath = getStringNotEmpty(configPath, exeDir + "/serverConfig") + "/";

	c.comMd.DbPath = c.comMd.ConfigPath + "data.db";
	c.comMd.WebConfigPath = c.comMd.ConfigPath + "config.xml";

	//load config
	c.cfgCtl = ConfigCtl{};
	c.cfgCtl.Load(c.comMd.WebConfigPath);

	c.comMd.Ip = c.cfgCtl.Md.Ip;
	c.comMd.Port = strconv.Itoa(c.cfgCtl.Md.Port);

	//init database
	dbServer := GetDbServer();
	dbServer.Init(c.comMd.DbPath);
	
	//init router
	c.router = MainRouter{};
	c.router.Init(c.comMd);

	//test

	//run http server
	url := "http://" + c.comMd.Ip + ":" + c.comMd.Port;
	fmt.Println("----------------------------------");
	fmt.Println("server start at: " + url);
	if(c.comMd.Ip == "0.0.0.0"){
		// listen all ip
		arrIp := c.findAllIp();
		mapIp := make(map[int]string)
		ch := make(chan ServerModel);
		for idx := range arrIp{
			// fmt.Println("ip: " + arrIp[idx]);
			mapIp[idx] = arrIp[idx];
			go c.runServerAsync(arrIp[idx], idx, ch);
		}

		go (func() {
			time.Sleep(1 * 500 * 1000000);
			strIp := "";
			arrListenIp := []string{};
			for idx := range mapIp {
				if strIp != "" {
					strIp += "|";
				}
				strIp += mapIp[idx];
				arrListenIp = append(arrListenIp, mapIp[idx]);
				// fmt.Println("ip: " + mapIp[idx]);
			}
			c.comMd.ArrIp = arrListenIp;
			fmt.Println("ips: " + strIp);
			fmt.Println("----------------------------------");
		})();

		// wait end
		for i := 0; i < 10; i++ {
			rst := <- ch;
			delete(mapIp, rst.Idx);
			// fmt.Println("aaa");
		}
		
	} else {
		arrListenIp := []string{};
		arrListenIp = append(arrListenIp, c.comMd.Ip);
		c.comMd.ArrIp = arrListenIp;

		// listen one ip
		c.runServer(c.comMd.Ip);
	}
}

type ServerModel struct {
	Idx		int;
	IsOk	bool;
}

func (c *ComSqliteServer) runServer(ip string){
	err := http.ListenAndServe(ip + ":" + c.comMd.Port, nil)
	if err != nil {
		// log.Fatal("ListenAndServe: ", err.Error());
		fmt.Println("ListenAndServe: ", err.Error());
	}
}

func (c *ComSqliteServer) runServerAsync(ip string, idx int, ch chan ServerModel){
	err := http.ListenAndServe(ip + ":" + c.comMd.Port, nil)
	if err != nil {
		// log.Fatal("ListenAndServe: ", err.Error());
		// fmt.Println(ip + " stopped: ", err.Error());
	}

	md := ServerModel{};
	md.Idx = idx;
	md.IsOk = (ch != nil);
	ch <- md;
}

func (c *ComSqliteServer) findAllIp() []string{
	arrIp := []string{}
	arrIp = append(arrIp, "localhost");
	// arrIp = append(arrIp, "127.0.0.1");

	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				// ip = v.IP;
				if !v.IP.IsLoopback() {
					if v.IP.To4() != nil {
						//Verify if IP is IPV4
						ip = v.IP
					}
				}
			// case *net.IPAddr:
			// 	ip = v.IP;
			}
			// process IP address
			str := ip.String();
			if str != "<nil>" {
				arrIp = append(arrIp, str);
			}
		}
	}

	return arrIp;
}

