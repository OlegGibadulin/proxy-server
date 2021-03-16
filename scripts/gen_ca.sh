#!/bin/sh

openssl genrsa -out ./scripts/ca.key 2048
openssl req -new -x509 -days 3650 -key ./scripts/ca.key -out ./scripts/ca.crt -subj "/CN=yngwie proxy CA"
openssl genrsa -out ./scripts/cert.key 2048
mkdir ./scripts/certs/
