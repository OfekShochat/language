package main

import (
	"fmt"
	"io/ioutil"
	"time"
)

func ReadFile(file_name string) string {
	dat, err := ioutil.ReadFile("file.g")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	for i := 0; i < len(dat); i++ {
		if dat[i] == '\r' {
			dat = append(dat[:i], dat[i+1:]...)
		}
	}
	return string(dat)
}

func main() {
	//var i32 = types.I32

	start := time.Now()

	s := ReadFile("file.g")

	out, err := Tokenize(s) //Tokenize("func f(i32 a, i32 b) -> i32 {\n  return 3 * 4\n}") //"func f(i32 a, i32 b) -> i32 {\n  return a * b\n}") // 1 + 2 - 3 + 1
	if err != nil {
		fmt.Println(err)
		return
	}
	/*for i := 0; i < len(out); i++ {
		fmt.Printf("%s '%s'\n", string(out[i].Type), out[i].Value)
	}*/
	_, functions := Parse(out)
	fmt.Println(functions)
	elapsed := time.Since(start)
	fmt.Println(float64(elapsed) / 1000000000)
	//fmt.Printf("%v\n", expressions[0])
	/*a, err := constant.NewIntFromString(i32, o[0].Params[1].Value)
	if err != nil {
		fmt.Println("poop")
	}
	b, err := constant.NewIntFromString(i32, o[0].Params[0].Value)
	if err != nil {
		fmt.Println("poop")
	}

	m := ir.NewModule()
	f := m.NewFunc("main", i32)
	f.NewBlock("")
	ff := f.Blocks[0]
	ff.NewRet(ff.NewAdd(a, b))
	fmt.Println(m.String())*/
}
