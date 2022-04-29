#! /bin/sh

openssl genrsa -aes256 -passout pass:gsahdg -out test.server.pass.key 4096
openssl rsa -passin pass:gsahdg -in test.server.pass.key -out test.server.key
openssl req -new -key test.server.key -out test.server.csr
openssl x509 -req -sha256 -days 365 -in test.server.csr -signkey test.server.key -out test.server.crt

rm -f test.server.pass.key test.server.csr
