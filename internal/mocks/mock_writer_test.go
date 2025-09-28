package mocks_test

import (
	"testing"

	"github.com/mr-filatik/go-password-keeper/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
	TestNewMockWriter
*/

func TestNewMockWriter(t *testing.T) {
	t.Parallel()

	mockWriter := mocks.NewMockWriter()

	require.NotNilf(t, mockWriter, "mockWriter should not be nil")
}

/*
	TestMockWriter_Write
*/

func TestMockWriter_Write(t *testing.T) {
	t.Parallel()

	mockWriter := mocks.NewMockWriter()

	number, err := mockWriter.Write([]byte("one"))

	require.NoError(t, err)
	require.NotEqual(t, 0, number)
}

/*
	TestMockWriter_MarkDataAsRead
*/

func TestMockWriter_MarkDataAsRead(t *testing.T) {
	t.Parallel()

	mockWriter := mocks.NewMockWriter()

	_, _ = mockWriter.Write([]byte("one"))
	_, _ = mockWriter.Write([]byte("two"))

	mockWriter.MarkDataAsRead()

	_, ok := mockWriter.GetUnreadedData()
	assert.False(t, ok, "there shouldn't be any unread data after MarkDataAsRead")
}

/*
	TestMockWriter_MarkDataAsRead
*/

func TestMockWriter_GetUnreadedData(t *testing.T) {
	t.Parallel()

	mockWriter := mocks.NewMockWriter()

	_, _ = mockWriter.Write([]byte("one"))
	_, _ = mockWriter.Write([]byte("two"))

	var data []byte

	var okData bool

	data, okData = mockWriter.GetUnreadedData()

	assert.True(t, okData, "data must be present")
	assert.Equal(t, []byte("one"), data)

	data, okData = mockWriter.GetUnreadedData()

	assert.True(t, okData, "data must be present")
	assert.Equal(t, []byte("two"), data)

	_, okData = mockWriter.GetUnreadedData()

	assert.False(t, okData, "data should not be present")
}
