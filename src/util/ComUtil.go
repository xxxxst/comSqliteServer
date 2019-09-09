package util

import (
	"io/ioutil"
	"regexp"
)

// 比较两个[]byte是否相同
func Equal(a []byte, b []byte) bool {
	// If one is nil, the other must also be nil.
    if (a == nil) != (b == nil) { 
        return false; 
    }

    if len(a) != len(b) {
        return false;
    }

    for i := range a {
        if a[i] != b[i] {
            return false;
        }
    }

    return true;
}

// func removeOneBom(bytes []byte, bom []byte) ([]byte, bool){
// 	isOk := Equal(bytes[:3], []byte{0xef, 0xbb, 0xbf});
// 	if(isOk){
// 		bytes = bytes[3:];
// 	}

// 	return bytes, isOk;
// }

// 去除文件bom
func removeBom(data []byte) []byte {
	arrBom := [][]byte {
		{0xef, 0xbb, 0xbf},			//utf-8
		{0xfe, 0xff},				//utf-16大端
		{0xff, 0xfe},				//utf-16小端
		{0x00, 0x00, 0xfe, 0xff},	//utf-32大端
		{0xff, 0xfe, 0x00, 0x00},	//utf-32小端
	}
	rst := data;
	isOk := false;
	for i:=0; i < len(arrBom); i++ {
		if len(data) < len(arrBom[i]) {
			continue;
		}
		
		isOk = Equal(data[:len(arrBom[i])], arrBom[i]);
		if(isOk){
			rst = data[len(arrBom[i]):];
			break;
		}
	}

	return rst;
}

// 读文件-去除bom
func ReadFile(path string) []byte {
	bytes,_ := ioutil.ReadFile(path);
	bytes = removeBom(bytes);
	return bytes;
}

// 读文件-去除bom
func ReadFileString(path string) string {
	bytes,_ := ioutil.ReadFile(path);
	bytes = removeBom(bytes);
	return string(bytes);
}

func Min(x, y int) int {
    if x <= y {
        return x
    }
    return y
}

func Max(x, y int) int {
    if x >= y {
        return x
    }
    return y
}

func SplitStr(text string, what string) []string {
	runeText := []rune(text);
	arrRst := []string{};

	arr := SplitRune(runeText, what);
	for _,val := range arr {
		arrRst = append(arrRst, string(val));
	}

	return arrRst;
}

func SplitRune(text []rune, what string) [][]rune {
	arrRst := [][]rune{};
	whatRunes := []rune(what);
	lenWaht := len(whatRunes);

	if len(text) == 0 {
		arrRst = append(arrRst, []rune{});
		return arrRst;
	}

	if lenWaht == 0 {
		arrRst = append(arrRst, text);
		return arrRst;
	}

	idx := 0;
    for i:=0; i<len(text); i++ {
        found := true
        for j := range whatRunes {
			if i+j >= len(text) {
				found = false;
				break;
			}
            if text[i+j] != whatRunes[j] {
                found = false
                break
            }
        }
        if found {
			// fmt.Println("aaa:", idx, i, text[idx:i]);
			arrRst = append(arrRst, text[idx:i]);
			i += lenWaht;
			idx = i;
        }
	}
	if idx < len(text) {
		arrRst = append(arrRst, text[idx:]);
	} else if idx == len(text) && idx != 0 {
		arrRst = append(arrRst, []rune{});
	}
    return arrRst;
}

func SearchStr(text string, what string) int {
	runeText := []rune(text);

	return SearchRune(runeText, what);
}

func SearchRune(text []rune, what string) int {
    whatRunes := []rune(what)

    for i := range text {
        found := true
        for j := range whatRunes {
			if i+j >= len(text) {
				return -1;
			}
            if text[i+j] != whatRunes[j] {
                found = false
                break
            }
        }
        if found {
            return i
        }
    }
    return -1
}

func FormatPath(path string) string {
	reg, _ := regexp.Compile("[\\/\\\\]+");
	path = reg.ReplaceAllString(path, "/");
	return path;
}