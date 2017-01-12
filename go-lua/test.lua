---- go-lua TEST FILE -----

GLOBAL_VAR = "this is a global var"

function test_go_string(fn, val)
  local str = fn(val)
  print("Type: " .. str)
end

function square(m)
  return m^2
end

function account_test() 
  local acc = Account.create(1000)
  
  print(acc)

  acc:withdrawl(100)
  print(acc:__tostring())
  print(acc:balance())
  print(Account.balance(acc))

  local acc2 = Account.create(900)
  print(acc2 == acc)
end

