#!/bin/bash
echo Waiting for pinner to start...
while ! nc -z ewm-das 3001; do
  sleep 1
done

echo -e "\nPinner is deployed"

echo -e "\nTesting endpoint /api/v1/cid..."
curl -F "filedata=@/root/.covalent/test/specimen.bin" http://ewm-das:3001/api/v1/cid

echo -e "\nTesting endpoint /api/v1/upload..."
curl -F "filedata=@/root/.covalent/test/specimen.bin" http://ewm-das:3001/api/v1/upload

echo -e "\nTesting endpoint /api/v1/download..."
curl http://ewm-das:3001/api/v1/download?cid=bafyreiaboifhabwe5opprpc2sfq2eab7n4csfksx2mywtg6sok4qcswblu

echo -e "\nTesting endpoint /upload..."
curl -F "filedata=@/root/.covalent/test/specimen.bin" http://ewm-das:3001/upload

echo -e "\nTesting endpoint /get..."
curl http://ewm-das:3001/get?cid=bafyreiaboifhabwe5opprpc2sfq2eab7n4csfksx2mywtg6sok4qcswblu

echo -e "\nTesting endpoint /cid..."
curl -F "filedata=@/root/.covalent/test/specimen.bin" http://ewm-das:3001/cid

exit 0
