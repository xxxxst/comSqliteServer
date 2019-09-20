// +build !debug

package tag

type TagDebug struct {
	Debug bool;
}

var GetTag = (func() (func() (*TagDebug)) {
	var ins *TagDebug;

	return func() (*TagDebug) {
		if(ins == nil) {
			ins = new(TagDebug);
			ins.Debug = false;
		}
		return ins;
	}
})();