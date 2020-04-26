#!/usr/bin/env sh

mkdir -p /etc/authentic/certificate
openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 \
    -subj "/C=UN/ST=MilkyWay/L=World/O=example/CN=localhost:8080" \
    -keyout /etc/authentic/certificate/tls.key \
    -out /etc/authentic/certificate/tls.crt


/authentic/dist/echo -bind localhost:8081 &
echoPid=$!

/authentic/dist/authentic proxy \
    --tls /etc/authentic/certificate \
    --target http://localhost:8081 \
    --bind localhost:8080 \
    --cookieName _session &
authenticPid=$!

sleep 1

cd /authentic/e2e
    /authentic/dist/check -workdir . -file main.lua
ec=$?

echo "Leaving"
sleep 1

kill ${echoPid}
kill ${authenticPid}
exit "${ec}"