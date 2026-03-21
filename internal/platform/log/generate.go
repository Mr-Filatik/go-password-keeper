//go:build generate
// +build generate

// Package log provides logging functionality.
package log

//go:generate mockgen -source=logger.go -destination=logger_mock_test.go -package=logging_test
