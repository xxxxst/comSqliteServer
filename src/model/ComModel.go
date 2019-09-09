package model

type ComModel struct {
	// ServerId		string;
	Version			string;
	
	ExePath			string;
	WebPath			string;
	DbPath			string;
	ConfigPath		string;
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
