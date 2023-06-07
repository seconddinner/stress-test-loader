#!/usr/bin/env bash

rm -rf client/cert
mkdir client/cert

cp cert/ca-key.pem client/cert/ca-key.pem
cp cert/ca-cert.pem client/cert/ca-cert.pem
cd cert

# 1. Generate web server's private key and certificate signing request (CSR)
echo "Generating server's private key and CSR"
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/O=seconddinner"

# 2. Use CA's private key to sign web server's CSR
echo "Signing server's CSR"
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem
openssl x509 -in server-cert.pem -noout -text

cd ../client/cert
# 3. Generate client's private key and certificate signing request (CSR)
echo "Generating client's private key and CSR"
openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-req.pem -subj "/O=seconddinner"

# 4. Use CA's private key to sign client's CSR
echo "Signing client's CSR"
openssl x509 -req -in client-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem
openssl x509 -in client-cert.pem -noout -text
