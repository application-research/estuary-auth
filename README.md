# Estuary Auth

Module for Estuary API authentication. It's to decouple the authorization from Estuary core to allow any API developers to
easily build authenticated APIs for Estuary.

This module can be imported as a module and added to the Estuary middleware func. 


![image](https://user-images.githubusercontent.com/4479171/179639246-2ae8c27c-fd9b-416f-8dda-be443f3d7526.png)


## Server (WIP)

## Import module
```
go get github.com/application-research/estuary-auth
```

## Initialize
```
//  initialize your database connection (estuary) - readonly

//  create the authorization middleware
authorizationServer := new(AuthorizationServer)
auth := authorizationServer.Init().SetDB(db).Connect()
```

## Use the middleware on your new Estuary API
```
//  add the authorization middleware to the Estuary middleware
//  PermLevelUpload = 1, PermLevelUser   = 2, PermLevelAdmin  = 10
e.GET("/metrics/", s.handleMetrics, auth.AuthRequired(authorization.PermLevelUser))
```
