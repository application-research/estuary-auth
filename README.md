# Estuary Auth Service API

[![Go](https://github.com/application-research/estuary-auth/actions/workflows/go.yml/badge.svg)](https://github.com/application-research/estuary-auth/actions/workflows/go.yml)

Estuary Auth Service is a service that provides authentication for Estuary. It's to decouple the authorization from Estuary core to allow any API developers to easily build authenticated APIs for Estuary.

![image](https://user-images.githubusercontent.com/4479171/179639246-2ae8c27c-fd9b-416f-8dda-be443f3d7526.png)


## Components
- Estuary Auth Core - the core authentication service library
- Estuary Auth Service API - the service rest api for Estuary Auth

## Running
create a .env with the following
```
DB_DSN=<DSN>
```

run the node
```
./estuary-auth-api
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


# View the usage docs here 
- To view the rest api go to [api](api).
- To view the  [core](core) library.
