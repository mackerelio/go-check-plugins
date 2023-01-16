openssl genrsa -out ca.key 2048

echo "generated ca.key"

openssl req -new -sha256 -key ca.key -out ca.csr -subj "/C=JP"

echo "generated ca.csr"

openssl x509 -in ca.csr -days 30 -req -signkey ca.key -sha256 -out ca.crt

echo "generated ca.crt"

openssl genrsa -out server.key 2048

echo "generated server.key"

openssl req -new -nodes -sha256 -key server.key -out server.csr -subj "/C=JP"

echo "generated server.csr"

openssl x509 -req -days 30 -in server.csr -sha256 -out server.crt -CA ca.crt -CAkey ca.key -CAcreateserial -extfile extfile.txt

echo "generated server.crt"

openssl genrsa -out client.key 2048

echo "generated client.key"

openssl req -new -nodes -sha256 -key client.key -out client.csr -subj "/C=JP"

echo "generated client.csr"

openssl x509 -req -days 30 -in client.csr -sha256 -out client.crt -CA ca.crt -CAkey ca.key -CAcreateserial -extfile extfile.txt

echo "generated client.crt"
