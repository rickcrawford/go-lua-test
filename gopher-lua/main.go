package main

import (
	"fmt"
	"log"

	lua "github.com/yuin/gopher-lua"
)

func runGlobalVar(L *lua.LState) {
	fmt.Printf("runGlobalVar, top stack: %d\n", L.GetTop())

	// Lua is a stack, so we get the object and check it's value off of it.
	val := L.GetGlobal("GLOBAL_VAR") // Load our global variable
	log.Printf("GLOBAL_VAR = %s\n", val)
}

func runInvalidVar(L *lua.LState) {
	fmt.Printf("runInvalidVar, top stack: %d\n", L.GetTop())

	// Lua is a stack, so we get the object and check it's value off of it.
	val := L.GetGlobal("ASDF_VAR") // Load our global variable
	if val.Type() != lua.LTNil {
		panic("this should not happen...")
	}
}

func runSquare(L *lua.LState) {
	fmt.Printf("runSquare, top stack: %d\n", L.GetTop())

	// call our square(int) function
	// if this function exists, we'll run it...
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("square"),
		NRet:    1,    // 1 return value
		Protect: true, // false return panic
	}, lua.LNumber(5)); err != nil {
		log.Fatalf("Error running square(5): %s", err)
	}

	ret := L.ToInt(-1) // returned value
	log.Printf("square(5) = %d\n", ret)

	L.Pop(1) // remove received value

}

func runGoTestFunc(L *lua.LState) {
	// show us the current index of the stack...
	fmt.Printf("runGoTestFunc, top stack: %d\n", L.GetTop())

	// Our simple function here will return true if it is a string, false if it isn't...
	ourSimpleFn := func(L *lua.LState) int {
		lv := L.Get(-1)
		switch lv.Type() {
		case lua.LTNumber:
			val, _ := lv.(lua.LNumber)
			log.Printf("fn(%d)\n", val)
			L.Push(lua.LString("int"))
		case lua.LTString:
			val, _ := lv.(lua.LString)
			log.Printf("fn(%s)\n", val)
			L.Push(lua.LString("string"))
		default:
			log.Printf("fn(%#v)\n", lv)
			L.Push(lua.LString("unknown"))
		}

		return 1 // number of return variables...
	}

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("test_go_string"),
		NRet:    0,    // 0 return value
		Protect: true, // false will panic
	}, L.NewClosure(ourSimpleFn), lua.LString("Hello, World!")); err != nil {
		log.Fatalf("Error running test_go_string: %s", err)
	}

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("test_go_string"),
		NRet:    0,    // 0 return value
		Protect: true, // false return panic
	}, L.NewClosure(ourSimpleFn), lua.LNumber(5)); err != nil {
		log.Fatalf("Error running test_go_string: %s", err)
	}

}

const accountName = "Account"

type Account struct {
	Balance int
}

func createAccount(L *lua.LState) int {
	log.Println("createAccount")

	account := &Account{Balance: L.CheckInt(1)}
	ud := L.NewUserData()
	ud.Value = account
	L.SetMetatable(ud, L.GetTypeMetatable(accountName))
	L.Push(ud)

	return 1
}

func accountToString(L *lua.LState) int {
	log.Println("accountToString")

	account := L.CheckUserData(1).Value.(*Account)
	L.Push(lua.LString(fmt.Sprintf("account(balance=%d)", account.Balance)))

	return 1
}

func accountEq(L *lua.LState) int {
	log.Println("accountEq")

	account1 := L.CheckUserData(1).Value.(*Account)
	account2 := L.CheckUserData(2).Value.(*Account)
	L.Push(lua.LBool(account1.Balance == account2.Balance))

	return 1
}

func accountWithdrawl(L *lua.LState) int {
	log.Println("accountWithdrawl")

	account := L.CheckUserData(1).Value.(*Account)
	account.Balance -= L.CheckInt(-1)

	return 0
}

func accountBalance(L *lua.LState) int {
	account := L.CheckUserData(1).Value.(*Account)
	L.Push(lua.LNumber(account.Balance))

	return 1
}

// http://stackoverflow.com/questions/34841773/golua-declaring-lua-class-with-defined-methods/34859174
func registerAccountType(L *lua.LState) {
	// Account = {}
	mt := L.NewTypeMetatable(accountName)

	L.SetGlobal(accountName, mt)
	// static attributes
	L.SetField(mt, "create", L.NewFunction(createAccount))
	L.SetField(mt, "balance", L.NewFunction(accountBalance))
	L.SetField(mt, "withdrawl", L.NewFunction(accountWithdrawl))
	L.SetField(mt, "__tostring", L.NewFunction(accountToString))
	L.SetField(mt, "__eq", L.NewFunction(accountEq))
	L.SetField(mt, "__index", mt)
}

func runMemberTest(L *lua.LState) {
	// register our type
	registerAccountType(L)

	fmt.Printf("runMemberTest, top stack: %d\n", L.GetTop())

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("account_test"),
		NRet:    0,    // 0 return value
		Protect: true, // false return panic
	}); err != nil {
		log.Fatalf("Error running account_test: %s", err)
	}

}

func main() {
	L := lua.NewState()
	L.OpenLibs()
	if err := L.DoFile("test.lua"); err != nil {
		log.Fatalf("Error loading file: %s", err)
	}

	runGlobalVar(L)

	runInvalidVar(L)

	runSquare(L)

	runGoTestFunc(L)

	runMemberTest(L)

	fmt.Printf("top stack: %d\n", L.GetTop())

}
