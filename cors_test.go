package cors

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSpec(t *testing.T) {
	cases := []struct {
		name       string
		cors       *Cors
		method     string
		reqHeaders map[string]string
		resHeaders map[string]string
		resCode    int
		resBody    string
	}{
		{
			name:   "PreflightRequest",
			cors:   &Cors{},
			method: http.MethodOptions,
			reqHeaders: map[string]string{
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "Authorization,x-Ccookie",
			},
			resHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Vary":                         "Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Methods": "GET",
				"Access-Control-Allow-Headers": "Authorization,x-Ccookie",
			},
			resCode: http.StatusNoContent,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, "http://example.com/foo", nil)

			for name, value := range tc.reqHeaders {
				req.Header.Add(name, value)
			}

			t.Run("Handler", func(t *testing.T) {
				res := httptest.NewRecorder()
				SetupRouter(tc.cors).ServeHTTP(res, req)
				assertResponse(t, res, tc.resCode, tc.resHeaders, tc.resBody)
			})
		})
	}
}

func SetupRouter(cors *Cors) *gin.Engine {
	router := gin.Default()

	router.Use(cors.GinHandler());

	router.GET("/foo", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Hello" : "World",
		})
	})

	return router
}

func assertResponse(
	t *testing.T,
	res *httptest.ResponseRecorder,
	httpCode int,
	headers map[string]string,
	body string,
) {
	a := assert.New(t)

	a.Equal(res.Code, httpCode)

	for k, v := range headers {
		a.Equal(res.Header().Get(k), v)
	}

	a.Equal(res.Body.String(), body)
}
