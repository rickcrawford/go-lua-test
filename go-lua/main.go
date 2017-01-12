package main

import (
	"fmt"
	"log"

	lua "github.com/Shopify/go-lua"
)

func runGlobalVar(L *lua.State) {
	fmt.Printf("runGlobalVar, top stack: %d\n", L.Top())

	// Lua is a stack, so we get the object and check it's value off of it.
	L.Global("GLOBAL_VAR") // Load our global variable
	// print the value
	if globalVarStr, ok := L.ToString(-1); ok {
		log.Printf("GLOBAL_VAR = %s\n", globalVarStr)
	}

	// clean up the stack...
	L.Pop(1)
}

func runInvalidVar(L *lua.State) {
	fmt.Printf("runInvalidVar, top stack: %d\n", L.Top())

	// Lua is a stack, so we get the object and check it's value off of it.
	L.Global("ASDF_VAR") // Load our global variable
	if _, ok := L.ToString(-1); ok {
		panic("this should not happen...")
	}

	// clean up stack... look another way!
	L.Remove(-1)
}

func runSquare(L *lua.State) {
	fmt.Printf("runSquare, top stack: %d\n", L.Top())

	// call our square(int) function
	L.Global("square")
	// if this function exists, we'll run it...
	if L.IsFunction(-1) {
		// put an element on the stack, in this case int(5)
		L.PushInteger(5)

		// the call function tells the stack that we are passing
		// one argument in and expecting one argument back
		if err := L.ProtectedCall(1, 1, 0); err != nil {
			log.Fatalf("Error running square: %s", err)
		}

		// we have a value on the stack we can get now, our result...
		if result, ok := L.ToInteger(-1); ok {
			log.Printf("square(5) = %d\n", result)
		}
	}
	// clean up the stack...
	L.Remove(-1)
}

func runGoTestFunc(L *lua.State) {
	// show us the current index of the stack...
	fmt.Printf("runGoTestFunc, top stack: %d\n", L.Top())

	// call our test go function...
	L.Global("test_go_string")
	// If this function exists, we'll run it
	if L.IsFunction(-1) {

		// Our simple function here will return true if it is a string, false if it isn't...
		ourSimpleFn := func(L *lua.State) int {
			switch {
			case L.IsNumber(-1):
				val, _ := L.ToInteger(-1)
				log.Printf("fn(%d)\n", val)
				L.PushString("int")
			case L.IsString(-1):
				val, _ := L.ToString(-1)
				log.Printf("fn(%s)\n", val)
				L.PushString("string")
			default:
				val := L.ToValue(-1)
				log.Printf("fn(%#v)\n", val)
				L.PushString("unknown")
			}

			return 1 // number of return variables...
		}

		// add a closure to the stack
		L.PushGoClosure(ourSimpleFn, 0)
		// add a string to the stack
		L.PushString("Hello, World!")

		// We're passing in 2 arguments, returning none...
		// This is equivalent to calling `test_go_string(ourSimpleFn, "Hello, World!")`
		if err := L.ProtectedCall(2, 0, 0); err != nil {
			log.Fatalf("Lua error: %#v\n", err)
		}

		// get our function back on the stack...
		L.Global("test_go_string")
		// add a closure to the stack
		L.PushGoClosure(ourSimpleFn, 0)
		// add a int to the stack
		L.PushInteger(123)

		// We're passing in 2 arguments, returning none...
		// This is equivalent to calling `test_go_string(ourSimpleFn, 123)`
		if err := L.ProtectedCall(2, 0, 0); err != nil {
			log.Fatalf("Lua error: %#v\n", err)
		}
	}
}

const accountName = "Account"

type Account struct {
	Balance int
}

func createAccount(L *lua.State) int {
	log.Println("createAccount")
	balance, _ := L.ToInteger(-1)
	account := &Account{Balance: balance}
	L.PushUserData(account)
	lua.SetMetaTableNamed(L, accountName)

	return 1
}

func accountToString(L *lua.State) int {
	log.Println("accountToString")
	account := lua.CheckUserData(L, 1, accountName).(*Account)
	L.PushString(fmt.Sprintf("account(balance=%d)", account.Balance))

	return 1
}

func accountEq(L *lua.State) int {
	log.Println("accountEq")
	account1 := lua.CheckUserData(L, 1, accountName).(*Account)
	account2 := lua.CheckUserData(L, 2, accountName).(*Account)
	L.PushBoolean(account1.Balance == account2.Balance)

	return 1
}

func accountWithdrawl(L *lua.State) int {
	log.Println("accountWithdrawl")
	account := L.ToUserData(1).(*Account)
	if amount, ok := L.ToInteger(-1); ok {
		account.Balance -= amount
	} else {
		lua.Errorf(L, "Invalid argument: %#v", L.ToValue(-1))
	}

	return 0
}

func accountBalance(L *lua.State) int {
	account := L.ToUserData(1).(*Account)
	L.PushInteger(account.Balance)

	return 1
}

// http://stackoverflow.com/questions/34841773/golua-declaring-lua-class-with-defined-methods/34859174
func registerAccountType(L *lua.State) {
	// Account = {}
	lua.NewMetaTable(L, accountName)
	// Account.__index = Account
	L.PushValue(-1)
	L.SetField(-2, "__index")

	lua.SetFunctions(L, []lua.RegistryFunction{
		{"create", createAccount},
		{"balance", accountBalance},
		{"withdrawl", accountWithdrawl},
		// reserved...
		{"__tostring", accountToString},
		{"__eq", accountEq},
	}, 0)

	// Add account to the global stack
	L.SetGlobal(accountName)
}

func runMemberTest(L *lua.State) {
	// register our type
	registerAccountType(L)

	fmt.Printf("runMemberTest, top stack: %d\n", L.Top())

	L.Global("account_test")
	if err := L.ProtectedCall(0, 0, 0); err != nil {
		log.Println(err)
	}
}

func main() {
	L := lua.NewState()
	lua.OpenLibraries(L)
	if err := lua.DoFile(L, "test.lua"); err != nil {
		log.Fatalf("Error loading file: %s", err)
	}

	runGlobalVar(L)

	runInvalidVar(L)

	runSquare(L)

	runGoTestFunc(L)

	runMemberTest(L)

	fmt.Printf("top stack: %d\n", L.Top())

}
