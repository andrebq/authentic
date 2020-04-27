local http = require('http')

local res, err = http.get('https://localhost:8080/auth/login', {})
if err then fatal(err) end

println(res.status_code)
println(res.body)