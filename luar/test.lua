---- LUAR TEST FILE -----

function pretty_json(val)
  -- call a global function defined that prints pretty json --
  local str = json.pretty(val)
  print("pretty json:\n" .. str)
end

function member_test()
  local p = person.new("rick")
  pretty_json(p)
  print(p)
  print("-----")
  print(p:Name())
  print("-----")
  print(p.Name())

end

function test_struct(obj) 
  print(obj.test())
  print(obj:test())
  print(obj.rick)
  print(obj['rick'])
  print(obj)
  print(obj.name)
  print(obj['name'])

  obj:add("1")
  obj:add("2")

end
