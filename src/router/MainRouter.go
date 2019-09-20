package router

import (
    "fmt"
    "io"
    "os"
	"net/http"
	"bufio"
    "path"
    "io/ioutil"
	"path/filepath"
	"encoding/json"
	"time"
	"strconv"
	"regexp"
	"os/exec"
	"bytes"
	// "strings"
	
	"github.com/gorilla/mux"
	
	. "server"
	. "model"
	util "util"
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
	c.Router.HandleFunc("/server/{project}/direct/query/{sql}", c.directQueryHandler).Methods("GET");
	c.Router.HandleFunc("/server/{project}/direct/query", c.directQueryPostHandler).Methods("POST");
	c.Router.HandleFunc("/server/{project}/direct/exec/{sql}", c.directExecHandler).Methods("GET");
	c.Router.HandleFunc("/server/{project}/direct/exec", c.directExecPostHandler).Methods("POST");
	c.Router.HandleFunc("/server/{project}/direct/insert/{sql}", c.directInsertHandler).Methods("GET");
	c.Router.HandleFunc("/server/{project}/direct/insert", c.directInsertPostHandler).Methods("POST");
	c.Router.HandleFunc("/server/{project}/query", c.queryPostHandler).Methods("POST");
	c.Router.HandleFunc("/server/{project}/exec", c.execPostHandler).Methods("POST");
	c.Router.HandleFunc("/server/{project}/insert", c.insertPostHandler).Methods("POST");
	c.Router.HandleFunc("/server/{project}/file/upload", c.uploadFileHandler).Methods("POST");
	c.Router.HandleFunc("/server/{project}/file/download", c.downloadFileHandler).Methods("POST");
	c.Router.HandleFunc("/server/{project}/file/delete", c.deleteFileHandler).Methods("POST");
	c.Router.HandleFunc("/server/{project}/directory/delete", c.deleteDirectoryHandler).Methods("POST");
	c.Router.HandleFunc("/server/{project}/directory/clear", c.clearDirectoryHandler).Methods("POST");
	c.Router.HandleFunc("/server/{project}/directory/list", c.directoryListHandler).Methods("POST");
	c.Router.HandleFunc("/server/{project}/directory/listAll", c.directoryListAllHandler).Methods("POST");
	c.Router.HandleFunc("/server/{project}/cmd", c.cmdHandler).Methods("POST");
	
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
	proj := params["project"];

	data := c.GetDbServer(proj).DirectQuery(sql);
	jsonData, _ := json.Marshal(data);
	writeGzipByte(w, r, jsonData);
}

//direct query Handler post
func (c *MainRouter) directQueryPostHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r);
	proj := params["project"];

	md := DirectQueryPostMd{};
	decoder := json.NewDecoder(r.Body);
	err := decoder.Decode(&md);
	if err != nil || md.Sql == "" {
		c.comErr(w, r, "参数错误");
		return;
	}

	data := c.GetDbServer(proj).DirectQuery(md.Sql);
	jsonData, _ := json.Marshal(data);
	writeGzipByte(w, r, jsonData);
	// io.WriteString(w, "abc");
}

//direct exec Handler
func (c *MainRouter) directExecHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r);
	sql := params["sql"];
	proj := params["project"];

	isOk := c.GetDbServer(proj).DirectExec(sql);
	rst := NewComRst();
	rst.Success = isOk;
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

//direct exec Handler post
func (c *MainRouter) directExecPostHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r);
	proj := params["project"];

	md := DirectQueryPostMd{};
	decoder := json.NewDecoder(r.Body);
	err := decoder.Decode(&md);
	if err != nil || md.Sql == "" {
		c.comErr(w, r, "参数错误");
		return;
	}

	isOk := c.GetDbServer(proj).DirectExec(md.Sql);
	rst := NewComRst();
	rst.Success = isOk;
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

//direct insert Handler
func (c *MainRouter) directInsertHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r);
	sql := params["sql"];
	proj := params["project"];

	id, err := c.GetDbServer(proj).DirectInsert(sql);
	rst := NewComRst();
	rst.Success = (err == nil);
	rst.Data = id;
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

//direct insert Handler post
func (c *MainRouter) directInsertPostHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r);
	proj := params["project"];

	md := DirectQueryPostMd{};
	decoder := json.NewDecoder(r.Body);
	err := decoder.Decode(&md);
	if err != nil || md.Sql == "" {
		c.comErr(w, r, "参数错误");
		return;
	}

	id, err := c.GetDbServer(proj).DirectInsert(md.Sql);
	rst := NewComRst();
	rst.Success = (err == nil);
	rst.Data = id;
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

//query Handler post
func (c *MainRouter) queryPostHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r);
	proj := params["project"];

	md := QueryPostMd{};
	decoder := json.NewDecoder(r.Body);
	err := decoder.Decode(&md);
	if err != nil || md.Sql == "" || md.Params == nil {
		c.comErr(w, r, "参数错误");
		return;
	}

	data := c.GetDbServer(proj).Query(md.Sql, md.Params);
	jsonData, _ := json.Marshal(data);
	writeGzipByte(w, r, jsonData);
	// io.WriteString(w, "abc");
}

//exec Handler post
func (c *MainRouter) execPostHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r);
	proj := params["project"];

	md := QueryPostMd{};
	decoder := json.NewDecoder(r.Body);
	err := decoder.Decode(&md);
	if err != nil || md.Sql == "" || md.Params == nil {
		c.comErr(w, r, "参数错误");
		return;
	}

	isOk := c.GetDbServer(proj).Exec(md.Sql, md.Params);
	rst := NewComRst();
	rst.Success = isOk;
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

//insert Handler post
func (c *MainRouter) insertPostHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r);
	proj := params["project"];

	md := QueryPostMd{};
	decoder := json.NewDecoder(r.Body);
	err := decoder.Decode(&md);
	if err != nil || md.Sql == "" || md.Params == nil {
		c.comErr(w, r, "参数错误");
		return;
	}

	id, err := c.GetDbServer(proj).Insert(md.Sql, md.Params);
	rst := NewComRst();
	rst.Success = (err == nil);
	rst.Data = id;
	if(err != nil) {
		rst.ErrInfo = err.Error();
	}
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

//upload file
func (c *MainRouter) uploadFileHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r);
	proj := params["project"];
	
	name := r.FormValue("path");
	rewrite := r.FormValue("rewrite");

	// strRename := r.FormValue("rename");
	rename,err := strconv.Atoi(r.FormValue("rename"));
	if err != nil {
		c.comErr(w, r, "参数错误");
		return;
	}
	
	file, handler, err := r.FormFile("file");
	if err != nil {
		c.comErr(w, r, "参数错误");
		return;
	}
	// fmt.Println(handler.Filename);

	if name == "" {
		name = handler.Filename;
	}

	fileName := name;
	if(rename == 1) {
		//获取文件后缀
		ext := path.Ext(handler.Filename);
		// fileName = util.FormatTime(time.Now(), "yyyy/MM/dd HH:mm:ss fff") + ext;
		fileName = util.FormatTime(time.Now(), "yyyyMMddHHmmssfff") + ext;
	}

	path := "";
	if(rewrite==""||rewrite=="1"||rewrite=="true") {
		// if(fileName[1] != ':') {
		// 	fileName = c.ComMd.ExePath + fileName;
		// }
		
		fileName = c.formatPath(fileName);
		if(fileName == "") {
			c.comErr(w, r, "参数错误");
			return;
		}

		path = c.GetProjPath(proj) + fileName;
	}

	dir := "";
	dir, fileName = filepath.Split(path);

	if(!util.DirectoryExists(dir)) {
		os.MkdirAll(dir, os.ModePerm);
	}

	var f *os.File = nil;
	if util.FileExists(path) {
        f, _ = os.OpenFile(path, os.O_WRONLY, 0666) 
        // fmt.Println("文件存在")
    } else {
        f, _ = os.Create(path)
        // fmt.Println("文件不存在")
	}
	
	if(f == nil){
		c.comErr(w, r, "失败");
		return;
	}

	writer := bufio.NewWriter(f);
	io.Copy(writer, file);
	writer.Flush();
	f.Close();

	rst := NewComRst();
	rst.Success = true;
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

//download file
func (c *MainRouter) downloadFileHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r);
	proj := params["project"];
	
	md := FileMd{};
	decoder := json.NewDecoder(r.Body);
	err := decoder.Decode(&md);
	if err != nil || md.Path == "" {
		io.WriteString(w, "");
		return;
	}

	fileName := md.Path;
	path := "";
	if(md.Rewrite=="" || md.Rewrite=="1" || md.Rewrite=="true") {
		fileName = c.formatPath(fileName);
		if(fileName == "") {
			io.WriteString(w, "");
			return;
		}
		path = c.GetProjPath(proj) + fileName;
	} else {
		if(fileName == "") {
			io.WriteString(w, "");
			return;
		}
	}
	
	if(!util.FileExists(path)) {
		io.WriteString(w, "");
		return;
	}

	_, fname := filepath.Split(path);

	w.Header().Add("Content-Type", "application/octet-stream");
	w.Header().Add("Content-Disposition", "attachment; filename=" + fname);

	// 读xlsx文件并发送
	f, _ := os.OpenFile(path, os.O_RDONLY, 0666);
	io.Copy(w, f);
	f.Close();
}

//run cmd
func (c *MainRouter) cmdHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r);
	proj := params["project"];
	
	md := CmdMd{};
	decoder := json.NewDecoder(r.Body);
	err := decoder.Decode(&md);
	if err != nil || md.Path == "" {
		io.WriteString(w, "");
		return;
	}

	fileName := md.Path;
	path := "";
	if(md.Rewrite=="" || md.Rewrite=="1" || md.Rewrite=="true") {
		fileName = c.formatPath(fileName);
		if(fileName == "") {
			io.WriteString(w, "");
			return;
		}
		path = c.GetProjPath(proj) + fileName;
	} else {
		if(fileName == "") {
			io.WriteString(w, "");
			return;
		}
	}
	
	if(!util.FileExists(path)) {
		io.WriteString(w, "");
		return;
	}

	// workDir, _ := os.Getwd();
	if(len(path) < 2 || path[1] != ':') {
		path = c.ComMd.WorkDir + "/" + path;
	}
	
	dir, _ := filepath.Split(path);

	// run cmd
	bout := bytes.NewBuffer(nil);
	berr := bytes.NewBuffer(nil);
	cmd := exec.Command(path, md.Args...);
	// cmd := exec.Command(path, md.Argument);
	cmd.Dir = dir;
	cmd.Stdout = bout;
	err = cmd.Run();
	str1 := bout.String();
	str1,_ = util.DecodeGbkStr(str1);
	// if str1 != "" && c.ComMd.Debug { fmt.Println(str1); }
	str2 := berr.String();
	str2,_ = util.DecodeGbkStr(str2);
	// if str2 != "" && c.ComMd.Debug { fmt.Println(str2); }
	str := str1;
	if(str1 != "") { str += "\r\n"; }
	str += str2;

	// back
	rst := NewComRst();
	rst.Success = (err == nil);
	rst.Data = str;
	if(err != nil) {
		rst.ErrInfo = err.Error();
	}
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

//delete file
func (c *MainRouter) deleteFileHandler(w http.ResponseWriter, r *http.Request){
	var err error = nil;
	path := c.logicGetFilePath(w, r);
	if(path == "") {
		return;
	}

	if(util.FileExists(path)) {
		err = os.Remove(path);
	}

	rst := NewComRst();
	rst.Success = (err == nil);
	if(err != nil) {
		rst.ErrInfo = err.Error();
	}
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

//delete directory
func (c *MainRouter) deleteDirectoryHandler(w http.ResponseWriter, r *http.Request){
	// var err error = nil;
	path,err := c.logicGetDirPath(w, r);
	if(err != nil || path == "") {
		return;
	}

	if(util.DirectoryExists(path)) {
		err = os.RemoveAll(path);
	}

	rst := NewComRst();
	rst.Success = (err == nil);
	if(err != nil) {
		rst.ErrInfo = err.Error();
	}
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

//clear directory
func (c *MainRouter) clearDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	// var err error = nil;
	path,err := c.logicGetDirPath(w, r);
	if(err != nil || path == "") {
		return;
	}

	if(util.DirectoryExists(path)) {
		rd, err1 := ioutil.ReadDir(path);
		if(err1 == nil) {
			for _, fi := range rd {
				if fi.IsDir() {
					err = os.RemoveAll(path+"/"+fi.Name());
				} else {
					err = os.Remove(path+"/"+fi.Name());
				}
			}
		}
	}

	rst := NewComRst();
	rst.Success = (err == nil);
	if(err != nil) {
		rst.ErrInfo = err.Error();
	}
	jsonData, _ := json.Marshal(rst);
	writeGzipByte(w, r, jsonData);
}

//directory list
func (c *MainRouter) directoryListHandler(w http.ResponseWriter, r *http.Request) {
	path,err := c.logicGetDirPath(w, r);
	if(err != nil) {
		return;
	}
	// if(path == "") {
	// 	return;
	// }
	if(path == "") {
		path = "./";
	}

	arr := []FileInfo{};

	if(util.DirectoryExists(path)) {
		rd, err1 := ioutil.ReadDir(path);
		if(err1 == nil) {
			for _, fi := range rd {
				md := FileInfo{};
				md.Name = fi.Name();
				md.IsDir = fi.IsDir();
				md.Children = []FileInfo{};
				arr = append(arr, md);
			}
		}
	}

	jsonData, _ := json.Marshal(arr);
	writeGzipByte(w, r, jsonData);
}

//directory list all
func (c *MainRouter) directoryListAllHandler(w http.ResponseWriter, r *http.Request) {
	path,err := c.logicGetDirPath(w, r);
	if(err != nil) {
		return;
	}
	if(path == "") {
		path = "./";
	}

	// fmt.Println(path);

	arr := []FileInfo{};
	if(util.DirectoryExists(path)) {
		arr = c.ergDirListHandler(path);
	}

	jsonData, _ := json.Marshal(arr);
	writeGzipByte(w, r, jsonData);
}

func (c *MainRouter) ergDirListHandler(path string) []FileInfo {
	arr := []FileInfo{};

	rd, err1 := ioutil.ReadDir(path);
	if(err1 == nil) {
		for _, fi := range rd {
			md := FileInfo{};
			md.Name = fi.Name();
			md.IsDir = fi.IsDir();
			md.Children = []FileInfo{};

			if(md.IsDir) {
				md.Children = c.ergDirListHandler(path + "/" + md.Name);
			}

			arr = append(arr, md);
		}
	}

	return arr;
}

func (c *MainRouter) logicGetFilePath(w http.ResponseWriter, r *http.Request) string {
	params := mux.Vars(r);
	proj := params["project"];
	
	md := FileMd{};
	decoder := json.NewDecoder(r.Body);
	err := decoder.Decode(&md);
	if err != nil || md.Path == "" {
		c.comErr(w, r, "参数错误");
		return "";
	}

	fileName := md.Path;
	if(md.Rewrite=="" || md.Rewrite=="1" || md.Rewrite=="true") {
		
	} else {
		if(fileName == "") {
			c.comErr(w, r, "参数错误");
			return "";
		}

		// if(fileName[1] != ':') {
		// 	fileName = c.ComMd.ExePath + fileName;
		// }

		return fileName;
	}

	fileName = c.formatPath(fileName);
	if(fileName == "") {
		c.comErr(w, r, "参数错误");
		return "";
	}
	return c.GetProjPath(proj) + fileName;
}

func (c *MainRouter) logicGetDirPath(w http.ResponseWriter, r *http.Request) (string,error) {
	params := mux.Vars(r);
	proj := params["project"];
	
	md := FileMd{};
	decoder := json.NewDecoder(r.Body);
	err := decoder.Decode(&md);
	if err != nil {
		c.comErr(w, r, "参数错误");
		return "",err;
	}

	fileName := md.Path;
	if(md.Rewrite=="" || md.Rewrite=="1" || md.Rewrite=="true") {

	} else {
		// if(fileName[1] != ':') {
		// 	fileName = c.ComMd.ExePath + fileName;
		// }
		return fileName,nil;
	}

	fileName = c.formatPath(fileName);
	return c.GetProjPath(proj) + fileName,nil;
}

func (c *MainRouter) formatPath(fileName string) string {
	reg1, _ := regexp.Compile("[ ]*[\\/\\\\]+[ ]*");
	fileName = reg1.ReplaceAllString(fileName, "/");

	reg2, _ := regexp.Compile("[^\\/]*[^\\.\\/][\\/]\\.\\.[\\/]");
	for ;; {
		if(!reg2.MatchString(fileName)) {
			break;
		}
		fileName = reg2.ReplaceAllString(fileName, "");
	}

	reg3, _ := regexp.Compile("([:*?<>|])|([^\\/]\\.\\/)");
	if(reg3.MatchString(fileName)) {
		return "";
	}

	return fileName;
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

func (c *MainRouter) GetDbServer(proj string) (*DbServer) {
	return GetDbServer(proj, c.ComMd.ConfigPath + proj + "/data.db");
}

func (c *MainRouter) GetProjPath(proj string) (string) {
	// return c.ComMd.ConfigPath + proj + "/data/";
	return c.ComMd.ConfigPath + proj + "/";
}