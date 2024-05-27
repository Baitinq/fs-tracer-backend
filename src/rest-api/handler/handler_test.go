package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Baitinq/fs-tracer-backend/lib"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestHandleGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := NewMockDB(ctrl)
	recorder := httptest.NewRecorder()

	handler := Handler{db: db}

	file := &lib.File{
		Absolute_path: "/tmp/file.txt",
	}
	db.EXPECT().GetLatestFileByPath(gomock.Any(), "/tmp/file.txt").Return(file, nil)

	handler.handleGet(recorder, httptest.NewRequest(http.MethodGet, "/file/%2ftmp%2Ffile.txt", nil))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, fmt.Sprintln("File: ", file), recorder.Body.String())
}
