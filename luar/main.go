package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
)

type Person struct {
	Name string `lua:"name"`
}

func (p *Person) F1(L *lua.State) int {
	val := L.ToString(-1)
	p.Name = val

	return 0
}

func (p *Person) F2(L *lua.State) int {
	L.PushString(p.Name)

	return 1
}

func (p Person) String() string {
	return fmt.Sprintf("Person{Name: %s}", p.Name)
}

func prettyJson(L *lua.State) {
	// load our global function to the stack
	L.GetGlobal("pretty_json")
	// add a struct to the stack
	luar.GoToLua(L, struct{ Name, Value string }{Name: "Rick", Value: "TEST"})
	// call our function
	L.Call(1, 0)
}

func memberTest(L *lua.State) {
	L.GetGlobal("member_test")
	L.Call(0, 0)
}

func addFunc(L *lua.State, key string, f func(*lua.State) int) {
	L.PushString(key)
	L.PushGoFunction(f)
	L.SetTable(-3)
}

func addString(L *lua.State, key, val string) {
	L.PushString(key)
	L.PushString(val)
	L.SetTable(-3)
}

func runStructTest(L *lua.State) {
	L.GetGlobal("test_struct")
	L.CreateTable(1, 0)
	// L.PushString("__index")
	// L.PushValue(-2)
	// L.SetTable(-3)

	addFunc(L, "test", func(L *lua.State) int {
		L.PushString("Asdf")
		return 1
	})

	addString(L, "rick", "test")
	if err := L.Call(1, 0); err != nil {
		log.Println(err)
	}
}

func runStructTest2(L *lua.State) {
	L.GetGlobal("test_struct")
	person := &Person{"RICK"}
	luar.GoToLua(L, person)

	addFunc(L, "test", func(L *lua.State) int {
		L.PushString("Asdf")
		return 1
	})

	str := []string{}

	addFunc(L, "add", func(L *lua.State) int {
		str = append(str, L.ToString(-1))
		return 0
	})

	if err := L.Call(1, 0); err != nil {
		log.Println(err)
	}

	fmt.Printf("%#v\n", str)
}

func runLuar(filename string) {
	log.Println("Running LUAC with LuaR bindings")

	L := lua.NewState()
	defer L.Close()
	L.OpenLibs()

	L.DoFile(filename)

	// // enable unicode otherwise things like string count will be invalid with multi-byte chars
	// unicode.GoLuaReplaceFuncs(L)

	// // setup constants and functions
	// luar.Register(L, "", luar.Map{
	// 	// default functions
	// 	"ipairs": luar.ProxyIpairs,
	// 	"pairs":  luar.ProxyPairs,
	// 	"type":   luar.ProxyType,
	// })

	// JSON pretty function...
	luar.Register(L, "json", luar.Map{
		"pretty": func(value interface{}) (string, error) {
			data, err := json.MarshalIndent(value, "", "\t")
			return string(data), err
		},
	})

	luar.Register(L, "person", luar.Map{
		"new": func(name string) *Person {
			return &Person{Name: name}
		},
	})

	// prettyJson(L)
	// memberTest(L)

	fmt.Println("--------------------")
	runStructTest(L)
	runStructTest2(L)

	fmt.Printf("top: %d\n", L.GetTop())

}

func main() {
	runLuar("test.lua")
}
