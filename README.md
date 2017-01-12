Go Lua Environment Tests
====================

Frustrated by a lack of documentation on the LuaJIT and Go I decided to do a quick project
to test out a couple of the popular libraries I came across. 



## [aarzilli/golua](https://github.com/aarzilli/golua)

aarzilli/golua is a cgo project wrapping the Lua 5.1 C api. You pair it with
[stevedonovan/luar](https://github.com/stevedonovan/luar) to make your life easier (use `v2` branch!)
since that project will provide proxy functions to go.

Challenges with this library are that it requires you to have the LuaJIT 5.1 runtime installed in your environment in order to run, and it uses `cgo` extentions to run it.

This API wraps Lua 5.1, and if you are using OSX, you probably have 5.2 installed. Run the following to install the lua 5.1 JIT using homebrew:

```bash
brew install lua51
```

At this point you should be able to kick off the lua5.1 env:

```bash
➜  go-lua-test lua5.1 -v
Lua 5.1.5  Copyright (C) 1994-2012 Lua.org, PUC-Rio
```

Next install the bindings:

```bash
go get -u github.com/aarzilli/golua/lua
```

#### Code Examples

LuaJIT is a stack, and we will push our commands onto it to run. This is not thread safe, but you can use pooling to instantiate more than one instance and then close them all when you are done.

Here's a simple example that initiaites state and then opens the basic libraries:

```go
// create a new VM
L := lua.NewState()
// this is a convinience function that opens all standard libraries
L.OpenLibs()
// close the VM
defer L.Close()
```

You may want to white-list the libaries you make available, and that is possible using the restricted library open functions. You can find them in this file: [https://github.com/aarzilli/golua/blob/master/lua/lua.go](https://github.com/aarzilli/golua/blob/master/lua/lua.go)

Example - open just base and math library:

```go
L := lua.NewState() // create a new VM
L.OpenBase() // open base library
L.OpenMath() // open math library
defer L.Close() // close the VM
```

LuaJIT is a stack that you push onto. In order to get a function or variable, you need to pop it into the stack.

In `luac/test.lua` there is a global variable `GLOBAL_VAR` that is a simple string. If I want to get the value of this string I would execute the following command:

```go
// Lua is a stack, so we get the object and check it's value off of it.
// Load our global variable to the stack
L.GetGlobal("GLOBAL_VAR")
// Get the value
globalVarStr := L.ToString(-1)
log.Printf("GLOBAL_VAR = %s\n", globalVarStr)
// Once we're done with our result variable, we're going to want to remove it from the stack.
// If we reuse this state, the stack could continue to grow!
L.Remove(-1)
```

If you want to call a function, you first need to load it in the stack, push the variables you want to the stack, and then call it.

Example from `luac/test.lua`:

```lua
-- return the square of m --
function square(m)
  return m^2
end
```

To call this function: 

```go
// load our function
L.GetGlobal("square")
// put an element on the stack, in this case int(5)
L.PushInteger(5)
// call the function (sorry ignoring the error...)
// the two integers represent the number of stack variables in, numuber of variables pushed
// to the stack out...  in this case, 1 in, 1 out square(5) returns 25
L.Call(1, 1)
// get our result
result := L.ToInteger(-1)
log.Printf("square(5)=%d\n", result)
// Cleanup...
L.Remove(-1) // or L.Pop(1)
```

Now this is a simple use case... But let's try a more exciting example where we pass in a user defined struct variable. I am creating a class `Account` similar to the one in this Lua [example](http://lua-users.org/wiki/SimpleLuaClasses
):

```lua
Account = {}
Account.__index = Account

function Account.create(balance)
   local acnt = {}             -- our new object
   setmetatable(acnt,Account)  -- make Account handle lookup
   acnt.balance = balance      -- initialize our object
   return acnt
end

function Account:withdraw(amount)
   self.balance = self.balance - amount
end

-- create and use an Account
acc = Account.create(1000)
acc:withdraw(100)

```
Now translate that class into a Go Struct:

```go
// Create the global name
const accountName = "Account" 

// Define the struct
type Account struct {
	Balance int64
}

// Now register our struct to the Lua state:
func register(L *lua.State) {
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
		
	// Add account to the global stack
	L.SetGlobal(accountName)
}

func createAccount(L *lua.State) int {
	// pop the integer value from the stack
	balance := L.ToInteger(-1)
	// get the pointer reference to our self object
	account := (*Account)(L.NewUserdata(uintptr(unsafe.Sizeof(Account{}))))
	L.LGetMetaTable(accountName)
	L.SetMetaTable(-2)
	// set the value
	account.Balance = int64(balance)
	// number of args returning
	return 1
}

func accountWithdrawl(L *lua.State) int {
	// get the self arg
	account := (*Account)(L.ToUserdata(1))
	// get the amount
	amount := L.ToInteger(-1)
	// set the balance
	account.Balance -= int64(amount)
	// number of args returning
	return 0
}

func accountBalance(L *lua.State) int {
	// get the self arg
	account := (*Account)(L.ToUserdata(1))
	// push the value
	L.PushInteger(account.Balance)
	// number of args returning
	return 1
}
```
Member variables in Lua are accessed using `act:balance()` This is functionally equivalent to `Account.balance(act)`, where `act == self`

If you look at `luac/main.go` there are more examples to play with.

There's a number of helper functions to help you to make sure you can test to see what the type of 

#### LuaR

[LuaR](https://github.com/stevedonovan/luar/tree/v2) is a helpful library that will wrap some of the terse stack code to make it easy to push/pop functions into the LuaJIT heap.  Keep in mind, `v2` branch
is the best one to use.

There is some sample code wrapping Go/Lua code in `luar/main.go`


## [shopify/go-lua](https://github.com/Shopify/go-lua)
This library has the Lua 5.2 VM implemented entirely in go! Sacraficing some performance for ultimate portability.

I have implemented the same functions in this embedded VM to show similarities to golua. From their own documentation portability comes at the sacrafice of performance - this is an order of magnitude slower than the C bindings.

## [yuin/gopher-lua](https://github.com/yuin/gopher-lua)
This is another all in one library. Since I went through the trouble for go-lua, I figured it couldn't hurt to have some more examples.



