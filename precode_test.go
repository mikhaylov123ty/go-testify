package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Выделим функцию для инициализации однотипных запросов
// count - количество кафе в запросе
// city - название города
func initTestRequest(count int, city string) *httptest.ResponseRecorder {
	target := fmt.Sprintf("/cafe?count=%d&city=%s", count, city)

	req := httptest.NewRequest("GET", target, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	return responseRecorder
}

// Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4

	resp := initTestRequest(7, "moscow")

	// проверка: код ответа 200, тело ответа не пустое
	require.Equal(t, http.StatusOK, resp.Code)
	require.NotEmpty(t, resp.Body)

	// разделение строки в массив и сравнение общего количества с количеством объектов в массиве
	cities := strings.Split(resp.Body.String(), ",")
	assert.Equal(t, len(cities), totalCount)
}

// Город, который передаётся в параметре city, не поддерживается. Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа
func TestMainHandlerWhenWrongCity(t *testing.T) {
	resp := initTestRequest(2, "moscoww")

	// проверка: код ответа 400, ошибка wrong city value в теле ответа
	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.NotEmpty(t, resp.Body)
	assert.Equal(t, "wrong city value", resp.Body.String())
}

// Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое
func TestMainHandlerWhenOK(t *testing.T) {
	resp := initTestRequest(2, "moscow")

	// проверка: код ответа 200, тело ответа не пустое
	require.Equal(t, http.StatusOK, resp.Code)
	require.NotEmpty(t, resp.Body)
}
