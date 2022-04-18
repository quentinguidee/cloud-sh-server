package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogin(testing *testing.T) {
	router := gin.New()

	Routes(router)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/login", nil)

	router.ServeHTTP(recorder, req)

	assert.Equal(testing, http.StatusOK, recorder.Code)
	assert.Equal(testing, `{"message":"Server is Running."}`, recorder.Body.String())
}
