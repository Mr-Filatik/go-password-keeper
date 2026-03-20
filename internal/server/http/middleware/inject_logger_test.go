package middleware_test

// import (
// 	"context"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
// 	"github.com/mr-filatik/go-password-keeper/internal/server/http/middleware"
// 	"github.com/stretchr/testify/assert"
// )

// type MockLogger struct {
// 	loggedMessages []string
// }

// //nolint:ireturn
// func (m *MockLogger) With(_ ...interface{}) logging.Logger {
// 	return m
// }

// func (m *MockLogger) Debug(msg string, _ ...interface{}) {
// 	m.loggedMessages = append(m.loggedMessages, "DEBUG: "+msg)
// }

// func (m *MockLogger) Info(msg string, _ ...interface{}) {
// 	m.loggedMessages = append(m.loggedMessages, "INFO: "+msg)
// }

// func (m *MockLogger) Warn(msg string, _ error, _ ...interface{}) {
// 	m.loggedMessages = append(m.loggedMessages, "WARN: "+msg)
// }

// func (m *MockLogger) Error(msg string, _ error, _ ...interface{}) {
// 	m.loggedMessages = append(m.loggedMessages, "ERROR: "+msg)
// }

// func (m *MockLogger) Fatal(msg string, _ error, _ ...interface{}) {
// 	m.loggedMessages = append(m.loggedMessages, "FATAL: "+msg)
// }

// func (m *MockLogger) Close() error {
// 	return nil
// }

// func TestInjectLogger(t *testing.T) {
// 	t.Parallel()

// 	tests := []struct {
// 		name          string
// 		logger        logging.Logger
// 		checkContext  func(*testing.T, context.Context)
// 		handlerCalled bool
// 	}{
// 		{
// 			name:   "logger successfully injected into context",
// 			logger: &MockLogger{},
// 			checkContext: func(t *testing.T, ctx context.Context) {
// 				t.Helper()

// 				// Проверяем, что логгер есть в контексте
// 				logger := logging.FromContext(ctx)
// 				assert.NotNil(t, logger, "logger should be present in context")

// 				// Проверяем, что это тот же логгер
// 				mockLogger, ok := logger.(*MockLogger)
// 				assert.True(t, ok, "logger should be of type *MockLogger")
// 				assert.NotNil(t, mockLogger)
// 			},
// 			handlerCalled: true,
// 		},
// 		{
// 			name:   "nil logger handling",
// 			logger: nil,
// 			checkContext: func(t *testing.T, ctx context.Context) {
// 				logger := logging.FromContext(ctx)
// 				assert.Nil(t, logger, "logger in context should be nil")
// 			},
// 			handlerCalled: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			// Создаем тестовый handler, который будет проверять наличие логгера в контексте
// 			handlerCalled := false
// 			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 				handlerCalled = true

// 				// Проверяем контекст
// 				if tt.checkContext != nil {
// 					tt.checkContext(t, r.Context())
// 				}

// 				w.WriteHeader(http.StatusOK)
// 			})

// 			// Применяем middleware
// 			middleware := middleware.InjectLogger(tt.logger)
// 			handler := middleware(next)

// 			// Создаем тестовый запрос
// 			req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
// 			recorder := httptest.NewRecorder()

// 			// Выполняем запрос
// 			handler.ServeHTTP(recorder, req)

// 			// Проверяем, что handler был вызван
// 			assert.Equal(t, tt.handlerCalled, handlerCalled, "handler called status mismatch")

// 			// Проверяем HTTP статус
// 			assert.Equal(t, http.StatusOK, recorder.Code, "response status should be 200 OK")
// 		})
// 	}
// }

// func TestInjectLogger_LoggerAccessInHandler(t *testing.T) {
// 	t.Parallel()

// 	// Создаем мок логгер
// 	mockLogger := &MockLogger{}

// 	// Создаем handler, который будет использовать логгер из контекста
// 	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Получаем логгер из контекста
// 		logger := logging.FromContext(r.Context())
// 		assert.NotNil(t, logger, "logger should be in context")

// 		// Используем логгер
// 		logger.Info("test message")
// 		logger.Debug("debug message")

// 		w.WriteHeader(http.StatusOK)
// 	})

// 	// Применяем middleware
// 	middleware := middleware.InjectLogger(mockLogger)
// 	handler := middleware(next)

// 	// Создаем и выполняем запрос
// 	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
// 	recorder := httptest.NewRecorder()

// 	handler.ServeHTTP(recorder, req)

// 	// Проверяем, что логгер был использован
// 	assert.Len(t, mockLogger.loggedMessages, 2, "should have 2 log messages")
// 	assert.Contains(t, mockLogger.loggedMessages[0], "test message")
// 	assert.Contains(t, mockLogger.loggedMessages[1], "debug message")
// }

// func TestInjectLogger_ContextPropagation(t *testing.T) {
// 	t.Parallel()

// 	// Создаем middleware цепочку для проверки сохранения контекста
// 	mockLogger := &MockLogger{}

// 	// Промежуточный middleware для добавления данных в контекст
// 	addValueMiddleware := func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			ctx := context.WithValue(r.Context(), "test-key", "test-value")
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}

// 	// Финальный handler для проверки
// 	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Проверяем, что оба значения присутствуют в контексте
// 		val := r.Context().Value("test-key")
// 		assert.Equal(t, "test-value", val, "custom context value should be preserved")

// 		logger := logging.FromContext(r.Context())
// 		assert.NotNil(t, logger, "logger should be in context")
// 		assert.Equal(t, mockLogger, logger, "logger should be the injected one")

// 		w.WriteHeader(http.StatusOK)
// 	})

// 	// Строим цепочку: сначала InjectLogger, потом addValueMiddleware
// 	handler := middleware.InjectLogger(mockLogger)(addValueMiddleware(next))

// 	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
// 	recorder := httptest.NewRecorder()

// 	handler.ServeHTTP(recorder, req)

// 	assert.Equal(t, http.StatusOK, recorder.Code)
// }
