# Estuary Auth Service API

This is the service rest api for Estuary Auth. It's to decouple the authorization from Estuary core to allow any API developers
to easily build authenticated APIs for Estuary.

## Running

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
./estuary-api
```

This opens up a port at 1313 by default

## Usage
### /check-api-key 
- URL: https://estuary-auth-api.onrender.com/check-api-key
- Method: POST
```
{
    "Token":"<token>"
}
```

### /check-user-api-key
- URL: https://estuary-auth-api.onrender.com/check-user-api-key
- Method: POST
```
{
    "Username":"alvinreyes",
    "Token":"<token>"
}
```
### /check-user-pass
- URL: https://estuary-auth-api.onrender.com/check-user-pass
- Method: POST
```
{
    "Username":"alvinreyes",
    "Password":"<password>"
}
```

# Remote endpoint
This service api is currently available here `https://estuary-auth-api.onrender.com/`
