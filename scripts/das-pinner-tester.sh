#!/bin/bash
echo Waiting for pinner to start...
while ! nc -z ewm-das 5080; do
  sleep 1
done

echo -e "\nPinner is deployed"

# Test the /api/v1/cid endpoint
echo -e "\nTesting endpoint /api/v1/cid..."
curl -s -X POST -F "filedata=@/root/.covalent/test/specimen-result.json" http://ewm-das:5080/api/v1/cid

# Test the /api/v1/upload endpoint
echo -e "\nTesting endpoint /api/v1/upload..."
curl -s -X POST -F "filedata=@/root/.covalent/test/specimen-result.json" http://ewm-das:5080/api/v1/upload

# Uncomment the following block to test the /api/v1/download endpoint
echo -e "\nTesting endpoint /api/v1/download..."
curl -s http://ewm-das:5080/api/v1/download?cid=bafyreif75ukvi7d6lsxca4lo325vgj5cwrmrfrsqjhqxf6mxjeyx3dua7y

# Test the /upload endpoint
echo -e "\nTesting endpoint /upload..."
curl -s -X POST -F "filedata=@/root/.covalent/test/specimen-result.json" http://ewm-das:5080/upload

# Uncomment the following block to test the /get endpoint
echo -e "\nTesting endpoint /get..."
curl -s http://ewm-das:5080/get?cid=bafyreif75ukvi7d6lsxca4lo325vgj5cwrmrfrsqjhqxf6mxjeyx3dua7y

# Test the /cid endpoint again
echo -e "\nTesting endpoint /cid..."
curl -s -X POST -F "filedata=@/root/.covalent/test/specimen-result.json" http://ewm-das:5080/cid

exit 0