package main

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/aarzilli/golua/lua"
)

// Lua documentation: http://www.lua.org/manual/5.1/manual.html
// Simple classes in Lua: http://lua-users.org/wiki/SimpleLuaClasses

func runGlobalVar(L *lua.State) {
	fmt.Printf("runGlobalVar, top stack: %d\n", L.GetTop())

	// Lua is a stack, so we get the object and check it's value off of it.
	L.GetGlobal("GLOBAL_VAR") // Load our global variable
	if L.IsString(-1) {
		// print the value
		globalVarStr := L.ToString(-1)
		log.Printf("GLOBAL_VAR = %s\n", globalVarStr)
	}
	// clean up the stack...
	L.Remove(-1)
}

func runInvalidVar(L *lua.State) {
	fmt.Printf("runInvalidVar, top stack: %d\n", L.GetTop())

	// Lua is a stack, so we get the object and check it's value off of it.
	L.GetGlobal("ASDF_VAR") // Load our global variable
	if !L.IsNil(-1) {
		panic("this should not happen...")
	}
	L.Remove(-1)
}

func runSquare(L *lua.State) {
	fmt.Printf("runSquare, top stack: %d\n", L.GetTop())

	// call our square(int) function
	L.GetGlobal("square")
	// if this function exists, we'll run it...
	if L.IsFunction(-1) {
		// put an element on the stack, in this case int(5)
		L.PushInteger(5)

		// the call function tells the stack that we are passing
		// one argument in and expecting one argument back
		if err := L.Call(1, 1); err != nil {
			// if err isn't nil, it's a special type - it will be *lua.LuaError
			if luaErr, ok := err.(*lua.LuaError); ok {
				log.Fatalf("Lua error: %#v\n", luaErr)
			}
		}

		// we have a value on the stack we can get now, our result...
		if L.IsNumber(-1) {
			log.Printf("square(5) = %d\n", L.ToInteger(-1))
		}
	}
	// clean up the stack...
	L.Remove(-1)
}

func runGoTestFunc(L *lua.State) {
	// show us the current index of the stack...
	fmt.Printf("runGoTestFunc, top stack: %d\n", L.GetTop())

	// call our test go function...
	L.GetGlobal("test_go_string")
	// If this function exists, we'll run it
	if L.IsFunction(-1) {

		// Our simple function here will return true if it is a string, false if it isn't...
		ourSimpleFn := func(L *lua.State) int {
			switch {
			case L.IsNumber(-1):
				log.Printf("fn(%d)\n", L.ToInteger(-1))
				L.PushString("int")
			case L.IsString(-1):
				log.Printf("fn(%s)\n", L.ToString(-1))
				L.PushString("string")
			default:
				log.Printf("fn(%#v)\n", L.ToGoStruct(-1))
				L.PushString("unknown")
			}

			return 1 // number of return variables...
		}

		// add a closure to the stack
		L.PushGoClosure(ourSimpleFn)
		// add a string to the stack
		L.PushString("Hello, World!")

		// We're passing in 2 arguments, returning none...
		// This is equivalent to calling `test_go_string(ourSimpleFn, "Hello, World!")`
		if err := L.Call(2, 0); err != nil {
			// if err isn't nil, it's a special type - it will be *lua.LuaError
			if luaErr, ok := err.(*lua.LuaError); ok {
				log.Fatalf("Lua error: %#v\n", luaErr)
			}
		}

		// get our function back on the stack...
		L.GetGlobal("test_go_string")
		// add a closure to the stack
		L.PushGoClosure(ourSimpleFn)
		// add a int to the stack
		L.PushInteger(123)

		// We're passing in 2 arguments, returning none...
		// This is equivalent to calling `test_go_string(ourSimpleFn, 123)`
		if err := L.Call(2, 0); err != nil {
			// if err isn't nil, it's a special type - it will be *lua.LuaError
			if luaErr, ok := err.(*lua.LuaError); ok {
				log.Fatalf("Lua error: %#v\n", luaErr)
			}
		}
	}
	// no need to call remove, nothing added to the stack by call...
}

const accountName = "Account"

type Account struct {
	Balance int64
}

func createAccount(L *lua.State) int {
	log.Println("createAccount")

	balance := L.ToInteger(-1)

	account := (*Account)(L.NewUserdata(uintptr(unsafe.Sizeof(Account{}))))
	L.LGetMetaTable(accountName)
	L.SetMetaTable(-2)

	account.Balance = int64(balance)

	return 1
}

func accountToString(L *lua.State) int {
	log.Println("accountToString")

	account := (*Account)(L.ToUserdata(1))
	L.PushString(fmt.Sprintf("account(balance=%d)", account.Balance))

	return 1
}

// If you are doing some global caching of objects, when the VM runs
// garbage collection this runs...
func accountGC(L *lua.State) int {
	log.Println("accountGC")

	account := (*Account)(L.ToUserdata(1))
	log.Printf("(account(balance = %d)).__gc()\n", account.Balance)

	return 0
}

func accountEq(L *lua.State) int {
	log.Println("accountEq")

	account1 := (*Account)(L.ToUserdata(1))
	account2 := (*Account)(L.ToUserdata(2))
	L.PushBoolean(account1.Balance == account2.Balance)

	return 1
}

func accountWithdrawl(L *lua.State) int {
	log.Println("accountWithdrawl")

	account := (*Account)(L.ToUserdata(1))
	amount := L.ToInteger(-1)
	account.Balance -= int64(amount)

	return 0
}

func accountBalance(L *lua.State) int {
	account := (*Account)(L.ToUserdata(1))
	L.PushInteger(account.Balance)

	return 1
}

// http://stackoverflow.com/questions/34841773/golua-declaring-lua-class-with-defined-methods/34859174
func registerAccountType(L *lua.State) {
	// Account = {}
	L.NewMetaTable(accountName)

	// Account.__index = Account
	L.LGetMetaTable(accountName) // load Account on stack
	L.SetField(-2, "__index")    // set index to Account

	// The 2 lines above are the same as the following:
	// L.PushString("__index")
	// L.PushValue(-2)
	// L.SetTable(-3)

	L.SetMetaMethod("create", createAccount)
	L.SetMetaMethod("balance", accountBalance)
	L.SetMetaMethod("withdrawl", accountWithdrawl)
	L.SetMetaMethod("__tostring", accountToString)
	L.SetMetaMethod("__gc", accountGC)
	L.SetMetaMethod("__eq", accountEq)

	// Add account to the global stack
	L.SetGlobal(accountName)
}

func runMemberTest(L *lua.State) {
	// register our type
	registerAccountType(L)

	fmt.Printf("runMemberTest, top stack: %d\n", L.GetTop())

	L.GetGlobal("account_test")
	if err := L.Call(0, 0); err != nil {
		log.Println(err)
	}
}

func runLuaC(filename string) {
	log.Println("Running LUAC bindings")

	// Initialize your state....
	L := lua.NewState() // create a new VM
	defer L.Close()     // close the VM

	// setup the libraries available and run our lua test file
	L.OpenBase() // open base library
	L.OpenMath() // open math library

	L.DoFile(filename) // panic if it doesn't compile...

	runGlobalVar(L)

	runInvalidVar(L)

	runSquare(L)

	runGoTestFunc(L)

	runMemberTest(L)

	fmt.Printf("top: %d\n", L.GetTop())

}

func main() {
	runLuaC("test.lua")
}
