# Go Image Service

This service is used to upload and take images.

## Allowed Endpoints and Methods

| Endpoint                                      | Method |
|-----------------------------------------------|-------|
| /images/:id[0-9]/:filename[a-z][A-Z].[a-z]{3} |POST|
| /images/:id[0-9]/:filename[a-z][A-Z].[a-z]{3} |GET|

## Installation & Run
### Download
```
    $ git clone https://github.com/ahmetberke/go-image-service
```

### Build & Run
```
    $ go mod download
```
```
    $ go build -o /image-service
    $ /image-service
```
```
    $ go run main.go
```


