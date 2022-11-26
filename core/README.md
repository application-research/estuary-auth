# Estuary Auth Core

Module for Estuary API authentication. It's to decouple the authorization from Estuary core to allow any API developers to
easily build authenticated APIs for Estuary.

This module can be imported as a module and added to the Estuary middleware func. 

## Import
```
go get github.com/application-research/estuary-auth/core
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