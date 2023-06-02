#!/usr/bin/env bash

rm -r cert
rm -r client/cert

mkdir cert
mkdir client/cert

cd cert
# Generate CA's private key and self-signed certificate
echo "Generating CA's private key and self-signed certificate"
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/O=seconddinner"
openssl x509 -in ca-cert.pem -noout -text