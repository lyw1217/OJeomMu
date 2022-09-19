package controller_test

import (
	"net/http"
	"net/http/httptest"
	"ojeommu/controller"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	controller.InitRoutes(router)
	controller.ServeStaticFiles(router)

	return router
}

func getRequest(t testing.TB, url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

var r *gin.Engine

func TestMain(m *testing.M) {
	r = SetUpRouter()

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestHomePage(t *testing.T) {
	req := getRequest(t, "/")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	//responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, http.StatusOK, w.Code)

	req = getRequest(t, "/index.html")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInfoPage(t *testing.T) {
	req := getRequest(t, "/info.html")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	//responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTestPage(t *testing.T) {
	req := getRequest(t, "/test.html")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	//responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchHandler(t *testing.T) {
	req := getRequest(t, "/sendToGo")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	//responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchBotHandler(t *testing.T) {
	req := getRequest(t, "/ojeommu")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	//responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWtImgHandler(t *testing.T) {
	req := getRequest(t, "/weather")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	//responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, http.StatusOK, w.Code)
}
