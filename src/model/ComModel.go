package model

type ComModel struct {
	Version			string;
	
	RootDir			string;
	// ExePath			string;
	// WorkDir			string;
	ConfigPath		string;

	// DbPath			string;
	WebPath			string;
	// DataPath		string;
	WebConfigPath	string;

	Ip				string;
	Port			string;

	ArrIp			[]string;
}

var GetComModel = (func() (func() (*ComModel)) {
	var instance *ComModel;

	return func() (*ComModel) {
		if(instance == nil) {
			instance = new(ComModel);
		}
		return instance;
	}
})();
