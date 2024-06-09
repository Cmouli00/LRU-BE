# LRU-BE


Use go version 1.18.7 


# If goversion 1.18.7 is not available, use GVM:

# For macOS:
  
  brew install gvm
  
  gvm install go1.18.7 
  
  gvm use go1.18.7


# To start:

go run .

Listening and serving HTTP on :8080


## Available APIs:

# To set a cache POST:

curl 'http://localhost:8080/cache/set' \
  -H 'Accept: application/json, text/plain, */*' \
  --data-raw '{"key":"1","value":"Test1","expiration":12000}'

# To get all GET:

curl 'http://localhost:8080/cache/getall' \
  -H 'Accept: application/json, text/plain, */*' 

# To get by Key GET:

curl 'http://localhost:8080/cache/get/1' \
  -H 'Accept: application/json, text/plain, */*' 

# To Delete by key DELETE:

curl 'http://localhost:8080/cache/delete/1' \
  -X 'DELETE' \
  -H 'Accept: application/json, text/plain, */*' 
