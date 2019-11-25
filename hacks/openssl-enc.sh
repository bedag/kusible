#!/bin/bash
# Encrypts a file using openssl in a way understood by kusible

if [ $# -ne 2 ]; then
  echo "Usage: ${0} <file> <password>"
fi

FILE="${1}"
PASS="${2}"

# PBKDF2 not supported because of https://github.com/Luzifer/go-openssl/tree/v3.1.0
openssl enc -aes-256-cbc -in "${FILE}" -k "${PASS}" -out "${FILE}.enc"