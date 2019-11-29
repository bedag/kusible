#!/bin/bash
# "Encrypts" a file using tar/7z in a way understood by kusible
# That this still has to be supported is purely for backwards compatibility

if [ $# -ne 2 ]; then
  echo "Usage: ${0} <file> <password>"
fi

FILE="${1}"
PASS="${2}"

tar -C "$(dirname "${FILE}")" -cpvf - "$(basename "${FILE}")"  | 7z a -bd -y -p"${PASS}" -mx="9" -mhe -t7z -si "${FILE}.enc.7z"