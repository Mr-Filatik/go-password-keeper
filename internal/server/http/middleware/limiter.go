// Package middleware provides functionality for HTTP middleware.
package middleware

// import (
// 	"context"
// 	"net"
// 	"net/http"
// 	"sync"

// 	"golang.org/x/time/rate"
// )

// // LimiterOpts — настройки rate limiter'а.
// type LimiterOpts struct {
// 	// Rate — средняя скорость, запросов в секунду (например, 5 rps).
// 	Rate rate.Limit

// 	// Burst — "всплеск": сколько запросов можно пропустить залпом,
// 	// пока не начнёт действовать Rate.
// 	Burst int

// 	// KeyFn — по какому ключу лимитировать (IP, userID, route и т.п.).
// 	// Если nil — используем IP клиента.
// 	KeyFn func(r *http.Request) string

// 	// Respond — как отвечать при превышении лимита.
// 	// Если nil — используем стандартный ответ 429 text/plain.
// 	Respond func(ctx context.Context, w http.ResponseWriter)
// }

// // Limiter — chi / net/http middleware.
// func Limiter(opts LimiterOpts) func(http.Handler) http.Handler {
// 	// Значения по умолчанию.
// 	if opts.Rate <= 0 {
// 		opts.Rate = 5 // 5 запросов в секунду
// 	}
// 	if opts.Burst <= 0 {
// 		opts.Burst = 10
// 	}
// 	if opts.KeyFn == nil {
// 		opts.KeyFn = clientIPKey
// 	}
// 	if opts.Respond == nil {
// 		opts.Respond = defaultTooManyRequests
// 	}

// 	// Мапа "ключ → rate.Limiter".
// 	var (
// 		mu       sync.Mutex
// 		limiters = make(map[string]*rate.Limiter)
// 	)

// 	// сюда вставить наш LRU кэш, через двусвязный список
// 	getLimiter := func(key string) *rate.Limiter {
// 		mu.Lock()
// 		defer mu.Unlock()

// 		lim, ok := limiters[key]
// 		if !ok {
// 			lim = rate.NewLimiter(opts.Rate, opts.Burst)
// 			limiters[key] = lim
// 		}
// 		return lim
// 	}

// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			key := opts.KeyFn(r)
// 			lim := getLimiter(key)

// 			if !lim.Allow() {
// 				// никакой паники — просто лог и ответ
// 				log := logctx.FromContext(r.Context())
// 				log.Warn("rate limit exceeded",
// 					"key", key,
// 					"remote_addr", r.RemoteAddr,
// 					"path", r.URL.Path,
// 				)

// 				opts.Respond(r.Context(), w)
// 				return
// 			}

// 			// всё ок — пускаем дальше по цепочке
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// // clientIPKey — простой ключ по IP (без X-Forwarded-For).
// func clientIPKey(r *http.Request) string {
// 	host, _, err := net.SplitHostPort(r.RemoteAddr)
// 	if err != nil {
// 		// на всякий случай, если формат неожиданный
// 		return r.RemoteAddr
// 	}
// 	return host
// }

// // defaultTooManyRequests — стандартный ответ 429.
// func defaultTooManyRequests(ctx context.Context, w http.ResponseWriter) {
// 	w.Header().Set("Retry-After", "1") // секунда, можно настроить
// 	w.WriteHeader(http.StatusTooManyRequests)
// 	_, _ = w.Write([]byte("too many requests\n"))
// }

// // func jsonTooManyRequests(ctx context.Context, w http.ResponseWriter) {
// // 	w.Header().Set("Content-Type", "application/json")
// // 	w.Header().Set("Retry-After", "1")
// // 	w.WriteHeader(http.StatusTooManyRequests)
// // 	_, _ = w.Write([]byte(`{"error":"too_many_requests"}`))
// // }
