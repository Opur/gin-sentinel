package gin_sentinel

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/gin-gonic/gin"
)

func TestLimiter(t *testing.T) {
	if err := sentinel.InitDefault(); err != nil {
		t.Error(err)
	}
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
		t.Error(err)
	}
	router := gin.New()
	userGroup := router.Group("/user")
	{
		userGroup.Use(Limiter("user", nil, sentinel.WithTrafficType(base.Inbound)))
		userGroup.GET("/", func(context *gin.Context) {
			context.Status(http.StatusOK)
		})
	}

	var (
		hs                  = httptest.NewServer(router)
		successCount uint64 = 0
		failedCount  uint64 = 0
		totalCount   uint64 = 100
		startWG             = sync.WaitGroup{}
		stopWG              = sync.WaitGroup{}
	)

	startWG.Add(int(totalCount))
	for i := 0; i < int(totalCount); i++ {
		stopWG.Add(1)
		go func() {
			startWG.Wait()
			defer stopWG.Done()

			resp, err := http.Get(hs.URL + "/user/")
			if err != nil {
				t.Fatal(err)
			}
			switch resp.StatusCode {
			case http.StatusOK:
				atomic.AddUint64(&successCount, 1)
			case http.StatusTooManyRequests:
				atomic.AddUint64(&failedCount, 1)
			default:
				t.Fatalf("unexcept status code: %d", resp.StatusCode)
			}
		}()
		startWG.Done()
	}
	stopWG.Wait()

	if successCount != 10 || successCount+failedCount != totalCount {
		t.Errorf("successCount: %d, failedCount: %d, totalCount: %d", successCount, failedCount, totalCount)
	}
}
