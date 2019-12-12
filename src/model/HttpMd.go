
package model

type ComRst struct {
	Success		bool				`json:"success"`;
	Data		interface{}			`json:"data"`;
	// InfoType	ErrType.ErrType		`json:"type"`;
	ErrInfo		string				`json:"errInfo"`;
}

func NewComRst() ComRst {
	rst := ComRst{};
	rst.Success = false;
	rst.Data = nil;
	rst.ErrInfo = "";

	return rst;
}

type DirectQueryPostMd struct {
	Sql			string		`json:"sql"`;
}

type QueryPostMd struct {
	Sql			string			`json:"sql"`;
	Params		[]interface{}	`json:"params"`;
}

type UploadFileMd struct {
	Rename		int	`json:"rename"`;
}

type FileMd struct {
	Path		string			`json:"path"`;
	Rewrite		string			`json:"rewrite"`;
}

type SaveStringFileMd struct {
	Path		string			`json:"path"`;
	Rewrite		string			`json:"rewrite"`;
	Rename		int				`json:"rename"`;
	Data		string			`json:"data"`;
}

type FileInfo struct {
	Name		string			`json:"name"`;
	IsDir		bool			`json:"isDir"`;
	Size		int64			`json:"size"`;
	ModifyTime	int64			`json:"modifyTime"`;
	// CreateTime	int64			`json:"createTime"`;
	Children	[]FileInfo		`json:"children"`;
}

type CmdMd struct {
	Path		string			`json:"path"`;
	Rewrite		string			`json:"rewrite"`;
	Args		[]string		`json:"args"`;
	// Argument	string			`json:"argument"`;
}