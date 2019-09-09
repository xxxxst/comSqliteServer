package control

import (
    "encoding/xml"
    "fmt"
    "io/ioutil"
	"os"
	
	. "model"
)

type XNodeConfig struct {
	XMLName			xml.Name	`xml:"ComSqliteServer"`
	Url				XNodeUrl	`xml:"url"`
}

type XNodeUrl struct {
	XMLName			xml.Name	`xml:"url"`
	
	IsHttps			bool		`xml:"https,attr"`
	Ip				string		`xml:"ip,attr"`
	Port			int			`xml:"port,attr"`
}

type ConfigCtl struct{
	Md WebConfigModel;
}

func (c *ConfigCtl) Load(path string){
	file, err := os.Open(path); // For read access.     
    if err != nil {
        fmt.Println("error: ", err);
        return
    }
	defer file.Close();
	
    data, err := ioutil.ReadAll(file);
    if err != nil {
        fmt.Println("error: ", err);
        return
	}
	
    v := XNodeConfig{}
    err = xml.Unmarshal(data, &v);
    if err != nil {
        fmt.Println("error: ", err);
        return
    }

	c.Md = WebConfigModel{};
	c.Md.IsHttps = v.Url.IsHttps;
	c.Md.Ip = v.Url.Ip;
	c.Md.Port = v.Url.Port;
	
	if(c.Md.Ip == ""){
		c.Md.Ip = "localhost";
	}
	if(c.Md.Port == 0){
		c.Md.Port = 9090;
	}
}
