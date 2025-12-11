//go:build generate
// +build generate

// Package logging provides logging functionality.
package logging

//go:generate mockgen -source=logger.go -destination=logger_mock_test.go -package=logging_test
