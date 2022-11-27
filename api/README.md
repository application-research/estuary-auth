# Estuary Auth Service API

This is the service rest api for Estuary Auth. It's to decouple the authorization from Estuary core to allow any API developers
to easily build authenticated APIs for Estuary.

## Running

```
go build -tags netgo -ldflags '-s -w' -o estuary-metrics-api
```
create a .env with the following
```
DB_NAME=
DB_HOST=
DB_USER=
DB_PASS=
DB_PORT=
```

run the node
```
./estuary-metrics-api
```

This opens up a port at 1313 by default

## Usage
### /check-api-key 
- URL: http://127.0.0.1:1313/check-api-key
- Method: POST
```
{
    "Token":"<token>"
}
```

### /check-user-api-key
- URL: http://127.0.0.1:1313/check-user-api-key
- Method: POST
```
{
    "Username":"alvinreyes",
    "Token":"<token>"
}
```
### /check-user-pass
- URL: http://127.0.0.1:1313/check-user-pass
- Method: POST
```
{
    "Username":"alvinreyes",
    "Password":"<password>"
}
```