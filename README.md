# Clean Arch in Go
[Read about Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## Build
```
  make
```

## Run tests
```
  make test
```

## Run API
```
  make run-api
```


## API requests 

### Add book

```
curl -X "POST" "http://localhost:8080/v1/book" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json' \
     -d $'{
  "title": "I Am Ozzy",
  "author": "Ozzy Osbourne",
  "pages": 294,
  "quantity":10
}'
```
### Search book

```
curl "http://localhost:8080/v1/book?title=ozzy" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json'
```

### Show books

```
curl "http://localhost:8080/v1/book" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json'
```