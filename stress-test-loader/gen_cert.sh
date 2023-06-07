#!/usr/bin/env bash

cd cert

# 1. Generate CA's certificate (given the CA's private key at ca-key.pem)
openssl req -new -x509 -key ca-key.pem -out ca-cert.pem -days 365 -subj "/O=seconddinner"

# 2. Generate web server's private key and certificate signing request (CSR)
echo "Generating server's private key and CSR"
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/O=seconddinner"

# 3. Use CA's private key to sign web server's CSR
echo "Signing server's CSR"
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem
openssl x509 -in server-cert.pem -noout -text

cd ..

rm -rf client/cert
mkdir client/cert

cp cert/ca-key.pem client/cert/ca-key.pem
cp cert/ca-cert.pem client/cert/ca-cert.pem

cd client/cert
# 4. Generate client's private key and certificate signing request (CSR)
echo "Generating client's private key and CSR"
openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-req.pem -subj "/O=seconddinner"

# 5. Use CA's private key to sign client's CSR
echo "Signing client's CSR"
openssl x509 -req -in client-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem
openssl x509 -in client-cert.pem -noout -text
