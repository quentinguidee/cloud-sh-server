package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAbout(testing *testing.T) {
	router := gin.New()
	router.GET("/", GetAbout)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	router.ServeHTTP(recorder, req)

	assert.Equal(testing, http.StatusOK, recorder.Code)
	assert.Equal(testing, `{"message":"Server is Running."}`, recorder.Body.String())
}
