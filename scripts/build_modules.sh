#!/bin/bash

echo "[!] Building Arsenal Modules"

cd arsenal || exit

go build -buildmode=c-shared -o arsenal.so arsenal.go

echo "[+] Modules build successfully"
