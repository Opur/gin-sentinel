package gin_sentinel

import (
	"net/http"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
)

// BlockHandler
type BlockHandler func(ctx *gin.Context, err *base.BlockError)

var defaultBlockHandler BlockHandler = func(ctx *gin.Context, err *base.BlockError) {
	ctx.String(http.StatusTooManyRequests, err.Error())
	ctx.Abort()
}

// SetDefaultBlockHandler set default BlockHandler.
func SetDefaultBlockHandler(handler BlockHandler) {
	if handler != nil {
		defaultBlockHandler = handler
	}
}

// Limiter return a new Limiter.
func Limiter(resource string, blockHandler BlockHandler, opt ...sentinel.EntryOption) gin.HandlerFunc {
	return func(context *gin.Context) {
		entry, err := sentinel.Entry(resource, opt...)
		if err != nil {
			if blockHandler == nil {
				blockHandler = defaultBlockHandler
			}
			blockHandler(context, err)
			return
		}
		// TODO: need exit opts here?
		defer entry.Exit()
		context.Next()
	}
}
