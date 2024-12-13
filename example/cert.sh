#!/bin/bash

CERTSDIR=certs
mkdir -p $CERTSDIR

function generate {

openssl req -x509                                     \
  -newkey rsa:4096                                    \
  -nodes                                              \
  -days 3650                                          \
  -keyout ${CERTSDIR}/${1}_ca_key.pem               \
  -out ${CERTSDIR}/${1}_ca_cert.pem                 \
  -subj /C=US/ST=CA/L=SVL/O=gRPC/CN=test-${1}_ca/   \

openssl genrsa -out ${CERTSDIR}/${1}_key.pem 4096

openssl req -new                                     \
  -key ${CERTSDIR}/${1}_key.pem                    \
  -out ${CERTSDIR}/${1}_csr.pem                    \
  -subj /C=US/ST=CA/L=SVL/O=gRPC/CN=test-${1}/    \

  -config openssl.cnf

openssl x509 -req                              \
  -in ${CERTSDIR}/${1}_csr.pem               \
  -CAkey ${CERTSDIR}/${1}_ca_key.pem         \
  -CA ${CERTSDIR}/${1}_ca_cert.pem           \
  -days 3650                                   \
  -set_serial 1000                             \
  -out ${CERTSDIR}/${1}_cert.pem             \
  -sha256                                      \
  -extfile openssl.cnf -extensions req_ext


openssl verify -verbose -CAfile ${CERTSDIR}/${1}_ca_cert.pem  ${CERTSDIR}/${1}_cert.pem

}

generate server
generate client