Lua Environment Test
====================

Frustrated by a lack of documentation on the LuaJIT and Go I decided to do a quick project
to test out a couple of the libraries I came across. 

[aarzilli/golua](https://github.com/aarzilli/golua)
----------------
In a past life this was the library I used to wrap the C bindings for LuaJIT in go. You pair it with
[stevedonovan/luar](https://github.com/stevedonovan/luar) to make your life easier (use `v2` branch!)
since that project will provide proxy functions to go.

** Installation **

Now this wraps Lua 5.1, and if you are using OSX, you probably have 5.2 installed. Run the following to install the lua 5.1 JIT using homebrew:

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

** Running ** 

LuaJIT is a stack that we will initialize and add our scripts to. This is not thread safe, but you can use pooling to instantiate more than one instance and then close them all when you are done.

```go
  // create a new VM
  L := lua.NewState()
  // this is a convinience function that opens all standard libraries
  L.OpenLibs()
  // close the VM
  defer L.Close()
```

You may want to white-list the libaries you make available, and that is possible using the restricted libary open functions.  You can find them in this file: [https://github.com/aarzilli/golua/blob/master/lua/lua.go](https://github.com/aarzilli/golua/blob/master/lua/lua.go)

Example - open just base and math library:

```go
  L := lua.NewState() // create a new VM
  L.OpenBase() // open base library
  L.OpenMath() // open math library
  defer L.Close() // close the VM
```

So LuaJIT is a stack that you push onto. In order to get a function or variable, you need to pop it into the stack.

In `lua/testluac.lua`, a global variable is defined: `GLOBAL_VAR` which is a simple string.  If I want to pull this value I can run the following commands:

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

Example from `lua/testluac.lua`

```lua
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
  L.Remove(-1)
```

Now this is a simple use case... But let's try a more exciting example where we pass in a user defined struct variable.




There's a number of helper functions to help you to make sure you can test to see what the type of 

** LuaR **

[LuaR](https://github.com/stevedonovan/luar/tree/v2) is a helpful library that will wrap some of the terse stack code to make it easy to push/pop functions into the LuaJIT heap.  Keep in mind, `v2` branch
is the best one to use in my humble opinion...



