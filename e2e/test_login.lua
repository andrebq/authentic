clear_cookies()

local http = require('http')
local url = require('url')

local res, err = http.get('https://localhost:8080/auth/login', {})
if err then fatal(err) end
if res.status_code ~= 200 then fatal("Unexpected status code.", res.status_code) end

local formBody = url.build_query_string({
    username= "nobody@example.com",
    password= "nobody"
})
println("formbody: ", formBody)

println("executing login")
res, err = http.post('https://localhost:8080/auth/login', {
    headers={
        ["Content-Type"]="application/x-www-form-urlencoded"},
    body= formBody})
if err then fatal(err) end
if res.status_code ~=200 then fatal("Unexpected status code.", res.status_code) end

local res, err = http.get('https://localhost:8080/hello/world')
if err then fatal(err) end
if res.status_code ~= 200 then fatal("Access should be allowed", res.status_code, " body: ", res.body) end
