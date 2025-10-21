//go:build docs

package server

// Ping godoc
// @Summary      Пинг сервиса
// @Description  Возвращает "pong" для проверки доступности сервиса.
// @Tags         health
// @Produce      plain
// @Success      200  {string}  string  "pong"
// @Failure      405  {string}  string  "method not allowed"
// @Failure      500  {string}  string  "internal server error"
// @Router       /ping [get]
