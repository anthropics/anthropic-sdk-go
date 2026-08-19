// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal/apiform"
	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// FileService contains methods and other services that help with interacting with
// the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewFileService] method instead.
type FileService struct {
	Options []option.RequestOption
}

// NewFileService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewFileService(opts ...option.RequestOption) (r FileService) {
	r = FileService{}
	r.Options = opts
	return
}

// List Files
func (r *FileService) List(ctx context.Context, query FileListParams, opts ...option.RequestOption) (res *pagination.PageCursor[FileMetadata], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/files"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// List Files
func (r *FileService) ListAutoPaging(ctx context.Context, query FileListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[FileMetadata] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

// Delete File
func (r *FileService) Delete(ctx context.Context, fileID string, opts ...option.RequestOption) (res *DeletedFile, err error) {
	opts = slices.Concat(r.Options, opts)
	if fileID == "" {
		err = errors.New("missing required file_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/files/%s", fileID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Download File
func (r *FileService) Download(ctx context.Context, fileID string, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/binary")}, opts...)
	if fileID == "" {
		err = errors.New("missing required file_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/files/%s/content", fileID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Get File Metadata
func (r *FileService) GetMetadata(ctx context.Context, fileID string, opts ...option.RequestOption) (res *FileMetadata, err error) {
	opts = slices.Concat(r.Options, opts)
	if fileID == "" {
		err = errors.New("missing required file_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/files/%s", fileID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Upload File
func (r *FileService) Upload(ctx context.Context, body FileUploadParams, opts ...option.RequestOption) (res *FileMetadata, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/files"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type DeletedFile struct {
	// ID of the deleted file.
	ID string `json:"id" api:"required"`
	// Deleted object type.
	//
	// For file deletion, this is always `"file_deleted"`.
	//
	// Any of "file_deleted".
	Type DeletedFileType `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DeletedFile) RawJSON() string { return r.JSON.raw }
func (r *DeletedFile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Deleted object type.
//
// For file deletion, this is always `"file_deleted"`.
type DeletedFileType string

const (
	DeletedFileTypeFileDeleted DeletedFileType = "file_deleted"
)

type FileMetadata struct {
	// Unique object identifier.
	//
	// The format and length of IDs may change over time.
	ID string `json:"id" api:"required"`
	// RFC 3339 datetime string representing when the file was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Original filename of the uploaded file.
	Filename string `json:"filename" api:"required"`
	// MIME type of the file.
	MimeType string `json:"mime_type" api:"required"`
	// Size of the file in bytes.
	SizeBytes int64 `json:"size_bytes" api:"required"`
	// Object type.
	//
	// For files, this is always `"file"`.
	Type constant.File `json:"type" default:"file"`
	// Whether the file can be downloaded.
	Downloadable bool `json:"downloadable"`
	// RFC 3339 datetime string representing when the file will expire and become
	// unavailable for download. Null if the file does not expire. For files uploaded
	// with `expires_in_seconds`, this is the upload time plus that value.
	ExpiresAt time.Time `json:"expires_at" api:"nullable" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		CreatedAt    respjson.Field
		Filename     respjson.Field
		MimeType     respjson.Field
		SizeBytes    respjson.Field
		Type         respjson.Field
		Downloadable respjson.Field
		ExpiresAt    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FileMetadata) RawJSON() string { return r.JSON.raw }
func (r *FileMetadata) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FileListParams struct {
	// Opaque page cursor returned in a prior list response's `next_page`. Prefixed
	// `page_`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `1000`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Restrict the result set to Files whose `id` is in this list. At most 100 entries
	// (after de-duplication). Mutually exclusive with `page` and `limit`. When
	// supplied, the response is always a single page (`next_page` is null). IDs that
	// do not resolve to a visible File — including deleted Files — are silently
	// omitted.
	IDs []string `query:"ids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [FileListParams]'s query parameters as `url.Values`.
func (r FileListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type FileUploadParams struct {
	// The file to upload
	File io.Reader `json:"file,omitzero" api:"required" format:"binary"`
	// Seconds from upload until the file expires and its bytes become permanently
	// unavailable. Must be between 3600 (one hour) and 7776000 (ninety days).
	ExpiresInSeconds param.Opt[int64] `json:"expires_in_seconds,omitzero"`
	paramObj
}

func (r FileUploadParams) MarshalMultipart() (data []byte, contentType string, err error) {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	err = apiform.MarshalRoot(r, writer)
	if err == nil {
		err = apiform.WriteExtras(writer, r.ExtraFields())
	}
	if err != nil {
		writer.Close()
		return nil, "", err
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), writer.FormDataContentType(), nil
}
