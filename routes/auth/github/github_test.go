package github

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	router := gin.New()
	err := godotenv.Load("../../../.env.test")
	if err != nil {
		panic(err)
	}
	LoadRoutes(router.Group("/auth"))

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/github/login", nil)

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, `{"url":"https://github.com/login/oauth/authorize?access_type=offline\u0026client_id=client_id\u0026redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Flogin\u0026response_type=code\u0026scope=all\u0026state=SH_CLOUD"}`, recorder.Body.String())
}
