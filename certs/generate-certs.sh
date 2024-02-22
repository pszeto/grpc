#/bin/sh
rm *.pem
rm *.srl
openssl req -x509 -nodes -days 73000 -newkey rsa:4096 -keyout ca-key.pem -out ca-cert.pem -subj "/CN=*.*"
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj  "/CN=*.*"
openssl x509 -req -in server-req.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.conf
openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-req.pem -subj "/CN=*.*"
openssl x509 -req -in client-req.pem -days 73000 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -extfile client-ext.conf
cp ca-cert.pem ../client
cp client-cert.pem ../client
cp client-key.pem ../client
cp ca-cert.pem ../server
cp server-cert.pem ../server
cp server-key.pem ../server
