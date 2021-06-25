package repository

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/pooladkhay/tinier-shortener-service/domain"
	"github.com/pooladkhay/tinier-shortener-service/helper/errs"
)

type Url interface {
	Create(url *domain.Url) *errs.Err
	GetByHash(hash string) (*domain.Url, *errs.Err)
	GetAll(user string) (*[]domain.Url, *errs.Err)
	Delete(hash, user string) *errs.Err
}

type url struct {
	client *gocql.Session
}

func NewUrl(s *gocql.Session) Url {
	return &url{client: s}
}

func (u *url) Create(url *domain.Url) *errs.Err {
	q := fmt.Sprintf("INSERT INTO urls (url, hash, user, is_private, created_at, expires_at, secret) VALUES (?, ?, ?, ?, ?, ?, ?) USING TTL %d;", url.ExpiresSecond)
	err := u.client.Query(
		q,
		url.Url,
		url.Hash,
		url.User,
		url.IsPrivate,
		url.CreatedAt,
		url.ExpiresAt,
		url.Secret,
	).Exec()
	if err != nil {
		return errs.NewBadRequestError(err.Error())
	}
	return nil
}

func (u *url) GetByHash(hash string) (*domain.Url, *errs.Err) {
	var result domain.Url
	var q = "SELECT url, hash, user, is_private, secret, created_at, expires_at FROM urls WHERE hash=?;"

	err := u.client.Query(q, hash).Scan(
		&result.Url,
		&result.Hash,
		&result.User,
		&result.IsPrivate,
		&result.Secret,
		&result.CreatedAt,
		&result.ExpiresAt,
	)
	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, errs.NewNotFoundError("invalid or expired url")
		}
		return nil, errs.NewInternalServerError(err.Error())
	}

	return &result, nil
}

func (u *url) GetAll(user string) (*[]domain.Url, *errs.Err) {
	q := fmt.Sprintf("SELECT url, hash, user, is_private, secret, created_at, expires_at FROM urls WHERE user = '%s' ALLOW FILTERING;", user)

	var urls []domain.Url

	scanner := u.client.Query(q).Iter().Scanner()
	for scanner.Next() {
		var url domain.Url
		err := scanner.Scan(&url.Url, &url.Hash, &url.User, &url.IsPrivate, &url.Secret, &url.CreatedAt, &url.ExpiresAt)
		if err != nil {
			if err == gocql.ErrNotFound {
				return nil, errs.NewNotFoundError("no url found for given hash")
			}
			return nil, errs.NewInternalServerError(err.Error())
		}
		urls = append(urls, url)
	}
	return &urls, nil
}

func (u *url) Delete(hash, user string) *errs.Err {
	q := "DELETE FROM urls WHERE hash=? AND user=? IF EXISTS;"

	err := u.client.Query(q, hash, user).Exec()
	if err != nil {
		if err == gocql.ErrNotFound {
			return errs.NewNotFoundError("no url found for given hash")
		}
		return errs.NewInternalServerError(err.Error())
	}
	return nil
}
