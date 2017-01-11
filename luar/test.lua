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
