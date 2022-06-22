package github

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogin(testing *testing.T) {
	router := gin.New()

	LoadRoutes(router.Group("/auth"))

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/github/login", nil)

	router.ServeHTTP(recorder, req)

	assert.Equal(testing, http.StatusOK, recorder.Code)
	assert.Equal(testing, `{"url":"https://github.com/login/oauth/authorize?access_type=offline\u0026client_id=\u0026redirect_uri=http%3A%2F%2Flocalhost%3A3000\u0026response_type=code\u0026scope=all"}`, recorder.Body.String())
}
