package main

import (
	"fmt"
	// "regexp"
)

func main() {
	comTest := ComTest{};
	comTest.Run();
}

type ComTest struct{
	// db *sqlx.DB;
}

func (c *ComTest) Run(){
	fmt.Println("test run");
	// Aaa();

	// reg := regexp.MustCompile(`.*\/([^\/]*)`);
	// match := reg.FindStringSubmatch(`default/doc/1314.md`);
	// fmt.Println("aaa:" + match[1]);

	// fmt.Println("1:", splitRune([]rune("abc"), "a"));
	// fmt.Println("2:", splitRune([]rune("abc"), "b"));
	// fmt.Println("3:", splitRune([]rune("abc"), "d"));
	// fmt.Println("4:", splitRune([]rune("abc"), ""));
	// fmt.Println("5:", splitRune([]rune(""), "a"));
	// fmt.Println("6:", splitRune([]rune("abc"), "c"));
	// fmt.Println("7:", splitRune([]rune("1abcad"), "a"));
	// fmt.Println("8:", splitRune([]rune("1abcabd"), "ab"));

	// fmt.Println("1:", splitStr("abc", "a"));
	// fmt.Println("2:", splitStr("abc", "b"));
	// fmt.Println("3:", splitStr("abc", "d"));
	// fmt.Println("4:", splitStr("abc", ""));
	// fmt.Println("5:", splitStr("", "a"));
	// fmt.Println("6:", splitStr("abc", "c"));
	// fmt.Println("7:", splitStr("1abcad", "a"));
	// fmt.Println("8:", splitStr("1abcabd", "ab"));
}

func Aaa(){
	fmt.Println("aaa");
	defer Bbb();
	fmt.Println("111");
}

func Bbb(){
	fmt.Println("bbb");
}