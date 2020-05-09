#!/usr/bin/env sh

redis-server >/dev/null 2>&1 &
redisPid=$!

mkdir -p /etc/authentic/certificate
openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 \
    -subj "/C=UN/ST=MilkyWay/L=World/O=example/CN=localhost:8080" \
    -keyout /etc/authentic/certificate/tls.key \
    -out /etc/authentic/certificate/tls.crt


/authentic/dist/echo -bind localhost:8081 &
echoPid=$!

export PROXY_TARGET="http://localhost:8081"
/authentic/dist/authentic proxy \
    --tls /etc/authentic/certificate \
    --bind localhost:8080 \
    --cookieName _session &
authenticPid=$!

sleep 1

cd /authentic/e2e
    /authentic/dist/check -workdir . -file main.lua
ec=$?

echo "Leaving"
sleep 1

/usr/bin/redis-cli <<EOF
keys *
EOF

kill ${echoPid}
kill ${authenticPid}
kill ${redisPid}
exit "${ec}"