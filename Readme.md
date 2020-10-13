# gin-sentinel

**gin-sentinel** is a rate limit middleware for [Gin](https://github.com/gin-gonic/gin)

## Installation

```shell script=
go get github.com/Opur/gin-sentinel
```

## Example

```go
package main

import (
    "net/http"
    
    ginSentinel "github.com/Opur/gin-sentinel"
    sentinel "github.com/alibaba/sentinel-golang/api"
    "github.com/alibaba/sentinel-golang/core/base"
    "github.com/alibaba/sentinel-golang/core/flow"
    "github.com/gin-gonic/gin"
)


func main() {
    // init sentinel first, see sentinel`s doc: https://github.com/alibaba/sentinel-golang/wiki
    if err := sentinel.InitDefault(); err != nil {
        panic(err)
 	}
    // add rules for sentinel, see sentinel`s doc: https://github.com/alibaba/sentinel-golang/wiki
 	_, err := flow.LoadRules([]*flow.Rule{
 		{
 			Resource:               "user",
 			MetricType:             flow.QPS,
 			Count:                  10,
 			TokenCalculateStrategy: flow.Direct,
 			ControlBehavior:        flow.Reject,
 		},
 	})
 	if err != nil {
 		panic(err)
 	}
    router := gin.New()
    userGroup := router.Group("/user")
    userGroup.Use(ginSentinel.Limiter("user", nil, sentinel.WithTrafficType(base.Inbound)))
    userGroup.GET("/", func(context *gin.Context) {
        context.Status(http.StatusOK)
    })
    router.Run(":8888")
}
```
