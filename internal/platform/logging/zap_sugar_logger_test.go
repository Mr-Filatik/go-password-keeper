package logging_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"errors"
// 	"strings"
// 	"testing"

// 	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// // helper: extracts the last non-empty line from the buffer and parses it as JSON.
// func parseLastJSONLine(t *testing.T, buf *bytes.Buffer) map[string]any {
// 	t.Helper()

// 	s := strings.TrimSpace(buf.String())
// 	require.NotEmpty(t, s, "expected some log output, got empty buffer")

// 	lines := strings.Split(s, "\n")
// 	last := strings.TrimSpace(lines[len(lines)-1])
// 	require.NotEmpty(t, last, "last log line is empty")

// 	var m map[string]any

// 	err := json.Unmarshal([]byte(last), &m)
// 	require.NoErrorf(t, err, "failed to unmarshal log line %q", last)

// 	return m
// }

// func TestNewZapSugarLogger_InitializesAndLogs(t *testing.T) {
// 	t.Parallel()

// 	var buf bytes.Buffer

// 	logger, err := logging.NewZapSugarLogger(logging.LevelInfo, &buf, logging.FormatJSON)
// 	require.NoError(t, err)
// 	require.NotNil(t, logger)

// 	lineMap := parseLastJSONLine(t, &buf)

// 	msg, ok := lineMap[logging.FieldBaseMessage].(string)
// 	require.True(t, ok, "expected %q field to be a string", logging.FieldBaseMessage)

// 	assert.Equal(t, "ZapSugar logger initialize is successful", msg)

// 	_, ok = lineMap[logging.FieldBaseLevel]
// 	assert.True(t, ok, "expected field %q (log level) to be present", logging.FieldBaseLevel)

// 	_, ok = lineMap[logging.FieldBaseData]
// 	assert.True(t, ok, "expected field %q (data) to be present", logging.FieldBaseData)
// }

// func TestZapSugarLogger_Debug_RespectsLogLevel(t *testing.T) {
// 	t.Parallel()

// 	t.Run("logs_on_debug_level", func(t *testing.T) {
// 		t.Parallel()

// 		var buf bytes.Buffer

// 		logger, err := logging.NewZapSugarLogger(logging.LevelDebug, &buf, logging.FormatJSON)
// 		require.NoError(t, err)

// 		buf.Reset()

// 		logger.Debug("debug-message", "foo", "bar")

// 		lineMap := parseLastJSONLine(t, &buf)

// 		msg, ok := lineMap[logging.FieldBaseMessage].(string)
// 		require.True(t, ok)
// 		assert.Equal(t, "debug-message", msg)

// 		data, ok := lineMap[logging.FieldBaseData].(map[string]any)
// 		require.True(t, ok, "expected %q to be object", logging.FieldBaseData)

// 		got, ok := data["foo"]
// 		assert.True(t, ok, "expected key %q in data", "foo")
// 		assert.Equal(t, "bar", got)
// 	})

// 	t.Run("does_not_log_on_info_level", func(t *testing.T) {
// 		t.Parallel()

// 		var buf bytes.Buffer

// 		logger, err := logging.NewZapSugarLogger(logging.LevelInfo, &buf, logging.FormatJSON)
// 		require.NoError(t, err)

// 		buf.Reset()

// 		logger.Debug("debug-message", "foo", "bar")

// 		assert.Equal(t, 0, buf.Len(), "expected no output for Debug with LevelInfo")
// 	})
// }

// func TestZapSugarLogger_Info_RespectsLogLevel(t *testing.T) {
// 	t.Parallel()

// 	t.Run("logs_when_level_allows", func(t *testing.T) {
// 		t.Parallel()

// 		var buf bytes.Buffer

// 		logger, err := logging.NewZapSugarLogger(logging.LevelInfo, &buf, logging.FormatJSON)
// 		require.NoError(t, err)

// 		buf.Reset()

// 		logger.Info("hello", "foo", "bar")

// 		lineMap := parseLastJSONLine(t, &buf)

// 		msg, ok := lineMap[logging.FieldBaseMessage].(string)
// 		require.True(t, ok)

// 		assert.Equal(t, "hello", msg)

// 		data, ok := lineMap[logging.FieldBaseData].(map[string]any)
// 		require.True(t, ok, "expected %q to be object", logging.FieldBaseData)

// 		got, ok := data["foo"]
// 		assert.True(t, ok, "expected key %q in data", "foo")
// 		assert.Equal(t, "bar", got)
// 	})

// 	t.Run("does_not_log_when_level_higher", func(t *testing.T) {
// 		t.Parallel()

// 		var buf bytes.Buffer

// 		logger, err := logging.NewZapSugarLogger(logging.LevelWarn, &buf, logging.FormatJSON)
// 		require.NoError(t, err)

// 		buf.Reset()

// 		logger.Info("hello", "foo", "bar")

// 		assert.Equal(t, 0, buf.Len(), "expected no output for Info with LevelWarn")
// 	})
// }

// func TestNewZapSugarLoggerWithFields_AddsCommonFields(t *testing.T) {
// 	t.Parallel()

// 	var buf bytes.Buffer

// 	loggerWithFields, err := logging.NewZapSugarLoggerWithFields(
// 		logging.LevelInfo,
// 		&buf,
// 		logging.FormatJSON,
// 		"app", "test-app",
// 	)
// 	require.NoError(t, err)
// 	require.NotNil(t, loggerWithFields)

// 	buf.Reset()

// 	loggerWithFields.Info("message")

// 	m := parseLastJSONLine(t, &buf)

// 	_, ok := m[logging.FieldBaseData]
// 	assert.True(t, ok, "expected %q to be present", logging.FieldBaseData)

// 	got, ok := m["app"]
// 	assert.True(t, ok, "expected top-level field %q to be present", "app")
// 	assert.Equal(t, "test-app", got)
// }

// func TestZapSugarLogger_With_AddsLabels(t *testing.T) {
// 	t.Parallel()

// 	var buf bytes.Buffer

// 	baseLogger, err := logging.NewZapSugarLogger(logging.LevelInfo, &buf, logging.FormatJSON)
// 	require.NoError(t, err)

// 	child := baseLogger.With("request_id", "req-123")
// 	require.NotNil(t, child)

// 	buf.Reset()

// 	child.Info("child-log")

// 	m := parseLastJSONLine(t, &buf)

// 	got, ok := m["request_id"]
// 	assert.True(t, ok, "expected request_id field to be present")
// 	assert.Equal(t, "req-123", got)
// }

// func TestZapSugarLogger_Warn_LogsErrorAndData(t *testing.T) {
// 	t.Parallel()

// 	var buf bytes.Buffer

// 	logger, err := logging.NewZapSugarLogger(logging.LevelWarn, &buf, logging.FormatJSON)
// 	require.NoError(t, err)

// 	buf.Reset()

// 	warnErr := errors.New("something went wrong")

// 	logger.Warn("warn-message", warnErr, "foo", "bar")

// 	assert.Contains(t, buf.String(), warnErr.Error(), "expected error to be present in log")
// }

// func TestZapSugarLogger_Error_LogsErrorAndData(t *testing.T) {
// 	t.Parallel()

// 	var buf bytes.Buffer

// 	logger, err := logging.NewZapSugarLogger(logging.LevelError, &buf, logging.FormatJSON)
// 	require.NoError(t, err)

// 	buf.Reset()

// 	errVal := errors.New("db connection failed")

// 	logger.Error("error-message", errVal, "foo", "bar")

// 	assert.Contains(t, buf.String(), errVal.Error(), "expected error to be present in log")
// }

// func TestZapSugarLogger_Close_DoesNotError(t *testing.T) {
// 	t.Parallel()

// 	var buf bytes.Buffer

// 	logger, err := logging.NewZapSugarLogger(logging.LevelInfo, &buf, logging.FormatJSON)
// 	require.NoError(t, err)

// 	err = logger.Close()
// 	assert.NoError(t, err)
// }
