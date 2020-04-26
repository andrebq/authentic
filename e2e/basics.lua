local http = require('http')

local res, err = http.get('http://localhost:8081/hello', {
    headers= {
        Cookie="_session=hello"
    }
})
if err then fatal(err) end

local res, err = http.get('https://localhost:8080/hello', {
    headers={
        Cookie="_session=hello"
    }
})
if err then fatal(err) end

if res.status_code ~= 200 then
    fatal("unexpected status code", res.status_code)
end