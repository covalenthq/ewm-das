#!/bin/bash
echo Waiting for pinner to start...
while ! nc -z ewm-das 3001; do
  sleep 1
done

echo -e "\nPinner is deployed"

echo -e "\nTesting endpoint /api/v1/cid..."
wget --quiet --method=POST --body-file=/root/.covalent/test/specimen.bin --header="Content-Type: multipart/form-data" -O - http://ewm-das:3001/api/v1/cid

echo -e "\nTesting endpoint /api/v1/upload..."
wget --quiet --method=POST --body-file=/root/.covalent/test/specimen.bin --header="Content-Type: multipart/form-data" -O - http://ewm-das:3001/api/v1/upload

echo -e "\nTesting endpoint /api/v1/download..."
wget --quiet -O - http://ewm-das:3001/api/v1/download?cid=bafyreiaboifhabwe5opprpc2sfq2eab7n4csfksx2mywtg6sok4qcswblu

echo -e "\nTesting endpoint /upload..."
wget --quiet --method=POST --body-file=/root/.covalent/test/specimen.bin --header="Content-Type: multipart/form-data" -O - http://ewm-das:3001/upload

echo -e "\nTesting endpoint /get..."
wget --quiet -O - http://ewm-das:3001/get?cid=bafyreiaboifhabwe5opprpc2sfq2eab7n4csfksx2mywtg6sok4qcswblu

echo -e "\nTesting endpoint /cid..."
wget --quiet --method=POST --body-file=/root/.covalent/test/specimen.bin --header="Content-Type: multipart/form-data" -O - http://ewm-das:3001/cid

exit 0