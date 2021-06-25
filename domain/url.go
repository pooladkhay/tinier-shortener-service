package domain

import (
	"time"

	"github.com/pooladkhay/tinier-shortener-service/helper/errs"
)

type Url struct {
	Url           string    `json:"url"`
	Hash          string    `json:"hash"`
	User          string    `json:"user"`
	IsPrivate     bool      `json:"is_private"`
	Secret        string    `json:"secret"`
	ExpiresSecond int       `json:"expires_second"`
	ExpiresAt     time.Time `json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
}

type ShortenRequest struct {
	Url           string `json:"url"`
	Hash          string `json:"hash"`
	IsPrivate     bool   `json:"is_private"`
	ExpiresSecond int    `json:"expires_second"`
	User          string
}
type ShortenResponse struct {
	Url       string    `json:"url"`
	Hash      string    `json:"hash"`
	User      string    `json:"user"`
	ShortUrl  string    `json:"short_url"`
	Secret    string    `json:"secret"`
	IsPrivate bool      `json:"is_private"`
	ExpiresAt time.Time `json:"expires_at"`
}

type GetByHashRequest struct {
	Hash   string `json:"hash"`
	Secret string `json:"secret"`
}
type GetByHashResponse struct {
	Url string `json:"url"`
}

type DeleteHashRequest struct {
	Hash   string `json:"hash"`
	Secret string `json:"secret"`
	User   string
}

func (sr *ShortenRequest) Validate() *errs.Err {
	if len(sr.Url) == 0 {
		return errs.NewBadRequestError("url is required")
	}
	if len(sr.Hash) > 0 && len(sr.Hash) < 6 {
		return errs.NewBadRequestError("hash should be at least 6 characters")
	}
	if !sr.IsPrivate {
		sr.IsPrivate = false
	}
	return nil
}
func (sr *GetByHashRequest) Validate() *errs.Err {
	if len(sr.Hash) == 0 {
		return errs.NewBadRequestError("hash is required")
	}
	if len(sr.Hash) > 0 && len(sr.Hash) < 6 {
		return errs.NewBadRequestError("hash should be at least 6 characters")
	}
	return nil
}
func (sr *DeleteHashRequest) Validate() *errs.Err {
	if len(sr.Hash) == 0 {
		return errs.NewBadRequestError("hash is required")
	}
	if len(sr.Hash) > 0 && len(sr.Hash) < 6 {
		return errs.NewBadRequestError("hash should be at least 6 characters")
	}
	return nil
}
