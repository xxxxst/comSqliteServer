
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
