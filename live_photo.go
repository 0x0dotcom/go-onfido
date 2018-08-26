package onfido

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"time"
)

// LivePhotoRequest represents a live photo request to Onfido API
type LivePhotoRequest struct {
	File               io.ReadSeeker
	ApplicantID        string
	AdvancedValidation *bool
}

// LivePhoto represents a live photo in Onfido API
type LivePhoto struct {
	ID           string     `json:"id,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	Href         string     `json:"href,omitempty"`
	DownloadHref string     `json:"download_href,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	FileType     string     `json:"file_type,omitempty"`
	FileSize     int        `json:"file_size,omitempty"`
}

// UploadLivePhoto uploads a live photo for the provided applicant.
// see https://documentation.onfido.com/#upload-live-photo
func (c *Client) UploadLivePhoto(ctx context.Context, dr LivePhotoRequest) (*LivePhoto, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := createFormFile(writer, "file", dr.File)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, dr.File); err != nil {
		return nil, err
	}
	if err := writer.WriteField("applicant_id", dr.ApplicantID); err != nil {
		return nil, err
	}
	av := "true"
	if dr.AdvancedValidation != nil && !*dr.AdvancedValidation {
		av = "false"
	}
	if err := writer.WriteField("advanced_validation", av); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := c.newRequest("POST", "/live_photos", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	var resp LivePhoto
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}
