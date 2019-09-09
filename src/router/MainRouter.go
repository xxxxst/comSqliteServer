package router

import (
    "fmt"
    // "io"
    "net/http"
	"encoding/json"
	
	"github.com/gorilla/mux"
	
	. "server"
	. "model"
)

type MainRouter struct{
	Router *mux.Router;
	// DbPath string;
	ComMd *ComModel;
}

func (c *MainRouter) Init(comMd *ComModel){
	c.ComMd = comMd;
	c.Router = mux.NewRouter();
	
	c.Router.HandleFunc("/server/{any:.*}", c.optionsHandler).Methods("OPTIONS");
	c.Router.HandleFunc("/server/direct/query/{sql}", c.directQueryHandler).Methods("GET");
	c.Router.HandleFunc("/server/direct/query", c.directQueryPostHandler).Methods("POST");
	c.Router.HandleFunc("/server/direct/exec/{sql}", c.directExecHandler).Methods("GET");
	c.Router.HandleFunc("/server/direct/exec", c.directExecPostHandler).Methods("POST");
	c.Router.HandleFunc("/server/direct/insert/{sql}", c.directInsertHandler).Methods("GET");
	c.Router.HandleFunc("/server/direct/insert", c.directInsertPostHandler).Methods("POST");
	c.Router.HandleFunc("/server/query", c.queryPostHandler).Methods("POST");
	c.Router.HandleFunc("/server/exec", c.execPostHandler).Methods("POST");
	c.Router.HandleFunc("/server/insert", c.insertPostHandler).Methods("POST");
	
	//static
	c.Router.HandleFunc("/{path:.*}", c.getStaticFileHandler);

	c.Router.Use(setOriginMiddleware);

	http.Handle("/", c.Router);
}

// 静态文件
func (c *MainRouter) getStaticFileHandler(w http.ResponseWriter, r *http.Request){
	getStaticFileHandler(w, r, c.ComMd.WebPath);
}

//direct query Handler
func (c *MainRouter) directQueryHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r);
	sql := params["sql"];

	data := GetDbServer().DirectQuery(sql);
	jsonData, _ := json.Marshal(data);
	writeGzipByte(w, r, jsonData);
}

//direct query Handler post
func (c *MainRouter) directQueryPostHandler(w http.ResponseWriter, r *http.Request){
	param := DirectQueryPostMd{};
	decoder := json.NewDecoder(r.Body);
	err := decoder.Decode(&param);
	if err != nil || param.Sql == "" {
		c.comErr(w, r, "参数错误");
		return;
	}

	data := GetDbServer().DirectQuery(param.Sql);
	jsonData, _ := json.Marshal(data);
	writeGzipByte(w, r, jsonData);
	// io.WriteString(w, "abc");
}

//direct exec Handler
func (c *MainRouter) directExecHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r);
	sql := params["sql"];

	isOk := GetDbServer().DirectExec(sql);
	rst := NewComRst();
	rst.Success = isOk;
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

//direct exec Handler post
func (c *MainRouter) directExecPostHandler(w http.ResponseWriter, r *http.Request){
	param := DirectQueryPostMd{};
	decoder := json.NewDecoder(r.Body);
	err := decoder.Decode(&param);
	if err != nil || param.Sql == "" {
		c.comErr(w, r, "参数错误");
		return;
	}

	isOk := GetDbServer().DirectExec(param.Sql);
	rst := NewComRst();
	rst.Success = isOk;
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

//direct insert Handler
func (c *MainRouter) directInsertHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r);
	sql := params["sql"];

	id, err := GetDbServer().DirectInsert(sql);
	rst := NewComRst();
	rst.Success = (err == nil);
	rst.Data = id;
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

//direct insert Handler post
func (c *MainRouter) directInsertPostHandler(w http.ResponseWriter, r *http.Request){
	param := DirectQueryPostMd{};
	decoder := json.NewDecoder(r.Body);
	err := decoder.Decode(&param);
	if err != nil || param.Sql == "" {
		c.comErr(w, r, "参数错误");
		return;
	}

	id, err := GetDbServer().DirectInsert(param.Sql);
	rst := NewComRst();
	rst.Success = (err == nil);
	rst.Data = id;
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

//query Handler post
func (c *MainRouter) queryPostHandler(w http.ResponseWriter, r *http.Request) {
	param := QueryPostMd{};
	decoder := json.NewDecoder(r.Body);
	err := decoder.Decode(&param);
	if err != nil || param.Sql == "" || param.Params == nil {
		c.comErr(w, r, "参数错误");
		return;
	}

	data := GetDbServer().Query(param.Sql, param.Params);
	jsonData, _ := json.Marshal(data);
	writeGzipByte(w, r, jsonData);
	// io.WriteString(w, "abc");
}

//exec Handler post
func (c *MainRouter) execPostHandler(w http.ResponseWriter, r *http.Request){
	param := QueryPostMd{};
	decoder := json.NewDecoder(r.Body);
	err := decoder.Decode(&param);
	if err != nil || param.Sql == "" || param.Params == nil {
		c.comErr(w, r, "参数错误");
		return;
	}

	isOk := GetDbServer().Exec(param.Sql, param.Params);
	rst := NewComRst();
	rst.Success = isOk;
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

//insert Handler post
func (c *MainRouter) insertPostHandler(w http.ResponseWriter, r *http.Request){
	param := QueryPostMd{};
	decoder := json.NewDecoder(r.Body);
	err := decoder.Decode(&param);
	if err != nil || param.Sql == "" || param.Params == nil {
		c.comErr(w, r, "参数错误");
		return;
	}

	id, err := GetDbServer().Insert(param.Sql, param.Params);
	rst := NewComRst();
	rst.Success = (err == nil);
	rst.Data = id;
	if(err != nil) {
		rst.ErrInfo = err.Error();
	}
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

// 返回错误信息
func (c *MainRouter) comErr(w http.ResponseWriter, r *http.Request, errInfo string) {
	rst := NewComRst();
	rst.ErrInfo = errInfo;
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

// 
func (c *MainRouter) comSuccess(w http.ResponseWriter, r *http.Request, data interface{}) {
	rst := NewComRst();
	rst.Data = data;
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

func (c *MainRouter) test(){
	fmt.Println("abc");
}
