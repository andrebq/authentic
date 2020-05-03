clear_cookies()
local http = require('http')

local res, err = http.get('https://localhost:8080/hello/world')
if err then fatal(err) end
if res.status_code ~= 401 then fatal("Access should be restricted", res.status_code, " body: ", res.body) end