package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Baitinq/fs-tracer-backend/lib"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestHandleGetFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := NewMockDB(ctrl)
	recorder := httptest.NewRecorder()

	handler := Handler{db: db}

	now := time.Now()
	file := &lib.File{
		Id:            "ID",
		User_id:       "USER_ID",
		Absolute_path: "/tmp/file.txt",
		Timestamp:     now,
		Contents:      "contents",
	}
	db.EXPECT().GetLatestFileByPath(gomock.Any(), "/tmp/file.txt", "USER_ID").Return(file, nil)

	handler.handleGet(recorder, httptest.NewRequest(http.MethodGet, "/file/?path=%2ftmp%2Ffile.txt", nil), "USER_ID")

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, strings.Join(strings.Fields(`
		{
			"Id": "ID",
			"User_id": "USER_ID",
			"Absolute_path": "/tmp/file.txt",
			"Contents": "contents",
			"Timestamp": "`+now.Format(time.RFC3339Nano)+`"
		}`), ""), recorder.Body.String())
}

func TestHandleGetRestoredFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := NewMockDB(ctrl)
	recorder := httptest.NewRecorder()

	handler := Handler{db: db}

	now := time.Now()
	file := &lib.File{
		Id:            "ID",
		User_id:       "USER_ID",
		Absolute_path: "/tmp/file.txt",
		Timestamp:     now,
		Contents:      "contents",
	}
	db.EXPECT().GetAndDeleteRestoredFiles(gomock.Any(), "USER_ID").Return(&[]lib.File{*file}, nil)

	handler.handleGet(recorder, httptest.NewRequest(http.MethodGet, "/restored-files/", nil), "USER_ID")

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, strings.Join(strings.Fields(`[
		{
			"Id": "ID",
			"User_id": "USER_ID",
			"Absolute_path": "/tmp/file.txt",
			"Contents": "contents",
			"Timestamp": "`+now.Format(time.RFC3339Nano)+`"
		}]`), ""), recorder.Body.String())
}
