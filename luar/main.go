package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
)

type Person struct {
	name string
}

func (p *Person) Name() string {
	return p.name
}

func (p Person) String() string {
	return fmt.Sprintf("Person{Name: %s}", p.name)
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
			return &Person{name: name}
		},
	})

	prettyJson(L)
	memberTest(L)

	fmt.Printf("top: %d\n", L.GetTop())

}

func main() {
	runLuar("test.lua")
}
