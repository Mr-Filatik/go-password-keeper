package app

type Option func(*App)

// // доп настройки
// func WithLaunchInterruption(laint bool) Option {
// 	return func(a *App) {
// 		a.launchInterruption = laint
// 	}
// }

// WithStopLaunchingOnError останавливает запуск других компонентов при получении ошибки.
// Пока действует на последовательный стартер.
func WithStopLaunchingOnError() Option {
	return func(a *App) {
		a.stopLaunchingOnError = true
	}
}
