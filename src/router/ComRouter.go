package router

import (
    "fmt"
    "io"
	"os"
    "path"
    "net/http"
	"net/url"
	"strings"
	"compress/gzip"
	"mime"
	"io/ioutil"
	
	// "github.com/gorilla/mux"
	
	// . "server"
	// . "model"
	// util "util"
)

func fileServer(dirPath string) http.Handler {
	return http.FileServer(http.Dir(dirPath));
}

func urlRedirect(converter func(string)(string), h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := converter(r.URL.Path);
		
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p
		h.ServeHTTP(w, r2)
	})
}

// static file
func getStaticFileHandler(w http.ResponseWriter, r *http.Request, webPath string){
	strPath := r.URL.Path;
	if strPath != "" && strPath[len(strPath)-1]=='/' {
		strPath = strPath + "index.html";
	}

	strPath = webPath + strPath;

	f, err := os.Stat(strPath);
	if err != nil || f.IsDir() {
		w.WriteHeader(404);
		writeGzipStr(w, r, "404 page not found");
		return;
	}

	sufix := path.Ext(strPath);
	conType := mime.TypeByExtension(sufix);
	if conType != "" {
        w.Header().Set("Content-Type", conType)
    }

	bytes,_ := ioutil.ReadFile(strPath);
	writeGzipByte(w, r, bytes);
}

func writeGzipByte(w http.ResponseWriter, r *http.Request, data []byte) {
	w.Write(data);

	// if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") && false {
	// 	w.Header().Set("Content-Encoding", "gzip");
	// 	gz := gzip.NewWriter(w);
	// 	defer gz.Close();
	// 	gz.Write(data);
	// } else {
	// 	// io.WriteString(w, string(data));
	// 	w.Write(data);
	// }
}

func writeGzipStr(w http.ResponseWriter, r *http.Request, data string) {
	// writeGzipByte(w, r, []byte(data));
	io.WriteString(w, string(data));
}

func writeGzipFile(w http.ResponseWriter, r *http.Request, f *os.File) {
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip");
		gz := gzip.NewWriter(w);
		defer gz.Close();
		io.Copy(gz, f);
	} else {
		io.Copy(w, f);
	}
}

func getFileContentType(path string) (string, error) {
	f, _ := os.Open(path);
	defer f.Close()
	
    // Only the first 512 bytes are used to sniff the content type.
    buffer := make([]byte, 512);

    _, err := f.Read(buffer);
    if err != nil {
        return "", err;
    }

    // Use the net/http package's handy DectectContentType function. Always returns a valid 
    // content-type by returning "application/octet-stream" if no others seemed to match.
    contentType := http.DetectContentType(buffer)

    return contentType, nil
}

func setOriginMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		str := r.Header.Get("Origin");
		w.Header().Add("Access-Control-Allow-Origin", str);
		w.Header().Add("Access-Control-Allow-Credentials", "true");
		
        next.ServeHTTP(w, r)
    })
}

func test(){
	fmt.Println("abc");
}

func (c *MainRouter) optionsHandler(w http.ResponseWriter, r *http.Request) {
	// setOrigin(w, r);

	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, X-Requested-With, Accept, Accept-Version, Content-Length, Content-MD5, Date, X-Api-Version, X-File-Name,Token,Cookie,authorization");
	w.Header().Add("Access-Control-Allow-Methods", "POST,GET,PUT,PATCH,DELETE,OPTIONS");
	io.WriteString(w, "");
}
