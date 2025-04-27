# Auction
Go expert auction project


### Environment variables

Set up environment variables on `cmd/auction/.env`
```bash
BATCH_INSERT_INTERVAL=20s
MAX_BATCH_SIZE=4
AUCTION_INTERVAL=1m
AUCTION_DURATION=1m
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=admin
MONGODB_URL=mongodb://admin:admin@mongodb:27017/auctions?authSource=admin
MONGODB_DB=auctions
```

### To build the executable
```bash
make build
```

### Run project
```bash
make run
```

### Stop project
```bash
make down
```

### Run test
```bash
make test
```

### Create an auction
```curl
curl --location 'localhost:8080/auction' \
--header 'Content-Type: application/json' \
--data '{
    "product_name": "Lord of The Rings",
    "category": "Fantasy",
    "description": "The history of middle earth",
    "condition": 0
}'
```
### Verify if auction was created
```curl
curl --location 'http://localhost:8080/auction?status=0' \
--header 'Content-Type: application/json'
```
### Verify if auction is closed
```curl
curl --location 'http://localhost:8080/auction?status=1' \
--header 'Content-Type: application/json'
```
