#!/bin/bash
echo Waiting for pinner to start...
while ! nc -z ewm-das 5080; do
  sleep 1
done

echo -e "\nPinner is deployed"

echo -e "\nTesting endpoint /api/v1/cid..."
wget --quiet --method=POST --body-file=/root/.covalent/test/specimen-result.json --header="Content-Type: multipart/form-data" -O - http://ewm-das:5080/api/v1/cid

echo -e "\nTesting endpoint /api/v1/upload..."
wget --quiet --method=POST --body-file=/root/.covalent/test/specimen-result.json --header="Content-Type: multipart/form-data" -O - http://ewm-das:5080/api/v1/upload

echo -e "\nTesting endpoint /api/v1/download..."
wget --quiet -O - http://ewm-das:5080/api/v1/download?cid=bafyreibo6rb2gvqi5srunoypym3tfzlbbj2yohmcbhpnrb47zexsugfeim

echo -e "\nTesting endpoint /upload..."
wget --quiet --method=POST --body-file=/root/.covalent/test/specimen-result.json --header="Content-Type: multipart/form-data" -O - http://ewm-das:5080/upload

echo -e "\nTesting endpoint /get..."
wget --quiet -O - http://ewm-das:5080/get?cid=bafyreibo6rb2gvqi5srunoypym3tfzlbbj2yohmcbhpnrb47zexsugfeim

echo -e "\nTesting endpoint /cid..."
wget --quiet --method=POST --body-file=/root/.covalent/test/specimen-result.json --header="Content-Type: multipart/form-data" -O - http://ewm-das:5080/cid

exit 0