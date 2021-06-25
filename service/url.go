package service

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/pooladkhay/tinier-shortener-service/domain"
	"github.com/pooladkhay/tinier-shortener-service/helper/errs"
	"github.com/pooladkhay/tinier-shortener-service/repository"
)

type Url interface {
	Shorten(*domain.ShortenRequest) (*domain.ShortenResponse, *errs.Err)
	GetByHash(hash, secret string) (*domain.GetByHashResponse, *errs.Err)
	GetAll(user string) (*[]domain.Url, *errs.Err)
	Delete(*domain.DeleteHashRequest) *errs.Err
}

type url struct {
	urlRepo   repository.Url
	cacheRepo repository.Cache
}

func NewUrl(u repository.Url, c repository.Cache) Url {
	return &url{urlRepo: u, cacheRepo: c}
}

func (u *url) Shorten(sr *domain.ShortenRequest) (*domain.ShortenResponse, *errs.Err) {
	newUrl := new(domain.Url)
	newUrl.Url = sr.Url
	newUrl.CreatedAt = time.Now().UTC()
	newUrl.User = sr.User

	// user wants a unique url (hash is provided in request)
	if len(sr.Hash) > 0 {
		if url, err := u.urlRepo.GetByHash(sr.Hash); url != nil && err == nil {
			return nil, errs.NewConflictError("hash already exists")
		}
		newUrl.Hash = sr.Hash
	} else {
		rand.Seed(time.Now().UnixNano())

		h := sha1.New()
		seed := fmt.Sprintf("%s+%d", sr.Url, time.Now().UTC().UnixNano())
		h.Write([]byte(seed))
		sha1Hash := hex.EncodeToString(h.Sum(nil))

		hashSlice := []string{}
		for i := 1; i < 7; i++ {
			index := rand.Intn(len(sha1Hash))
			hashSlice = append(hashSlice, string(rune(sha1Hash[index])))
		}
		newUrl.Hash = strings.Join(hashSlice, "")
	}

	// user wants a private url (is_private is true)
	if sr.IsPrivate {
		newUrl.IsPrivate = true
		seed := fmt.Sprintf("%s+%d", sr.Url, time.Now().UTC().UnixNano())
		hash := md5.Sum([]byte(seed))
		newUrl.Secret = hex.EncodeToString(hash[:])[:9]
	} else {
		newUrl.IsPrivate = false
		seed := fmt.Sprintf("%s+%d", sr.Url, time.Now().UTC().UnixNano())
		hash := md5.Sum([]byte(seed))
		newUrl.Secret = hex.EncodeToString(hash[:])[:9]
	}

	// user has provided expiration time in second
	if sr.ExpiresSecond > 0 {
		newUrl.ExpiresAt = time.Now().UTC().Add(time.Second * time.Duration(sr.ExpiresSecond))
		newUrl.ExpiresSecond = sr.ExpiresSecond
	} else {
		dayInSecond := 24 * 60 * 60
		newUrl.ExpiresAt = time.Now().UTC().Add(time.Second * time.Duration(dayInSecond))
		newUrl.ExpiresSecond = dayInSecond
	}

	// finally, save new url to db
	if err := u.urlRepo.Create(newUrl); err != nil {
		return nil, err
	}

	if newUrl.IsPrivate {
		u.cacheRepo.Cache(newUrl.Hash, fmt.Sprintf("%s|||%s", newUrl.Secret, newUrl.Url), newUrl.ExpiresSecond)
	} else {
		u.cacheRepo.Cache(newUrl.Hash, newUrl.Url, newUrl.ExpiresSecond)
	}

	// response to user
	resp := new(domain.ShortenResponse)
	resp.Url = newUrl.Url
	resp.Hash = newUrl.Hash
	resp.User = newUrl.User
	resp.Secret = newUrl.Secret
	resp.IsPrivate = newUrl.IsPrivate
	resp.ExpiresAt = newUrl.ExpiresAt
	resp.ShortUrl = fmt.Sprintf("%s/%s", os.Getenv("MAIN_URL"), newUrl.Hash)
	if newUrl.IsPrivate {
		resp.ShortUrl = fmt.Sprintf("%s/%s?secret=%s", os.Getenv("MAIN_URL"), newUrl.Hash, newUrl.Secret)
	}

	return resp, nil
}

func (u *url) GetByHash(hash, secret string) (*domain.GetByHashResponse, *errs.Err) {
	resp := new(domain.GetByHashResponse)

	cUrl, _ := u.cacheRepo.GetCache(hash)
	if cUrl != nil {
		if strings.Contains(*cUrl, "|||") {
			s := strings.Split(*cUrl, "|||")
			if s[0] == secret {
				resp.Url = s[1]
				return resp, nil
			} else {
				return nil, errs.NewUnauthorizedError("url is private")
			}
		}
		resp.Url = *cUrl
		return resp, nil
	}

	url, err := u.urlRepo.GetByHash(hash)
	if err != nil {
		return nil, err
	}
	if url.IsPrivate {
		if url.Secret == secret {
			resp.Url = url.Url
			return resp, nil
		} else {
			return nil, errs.NewUnauthorizedError("url is private")
		}
	}
	resp.Url = url.Url
	return resp, nil
}

func (u *url) GetAll(user string) (*[]domain.Url, *errs.Err) {
	urls, err := u.urlRepo.GetAll(user)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func (u *url) Delete(d *domain.DeleteHashRequest) *errs.Err {

	url, err := u.urlRepo.GetByHash(d.Hash)
	if err != nil {
		return err
	}

	if d.User != "" {
		if d.User == url.User {
			return u.urlRepo.Delete(d.Hash, d.User)
		}
	} else {
		if d.Secret == url.Secret {
			return u.urlRepo.Delete(d.Hash, d.User)
		}
	}

	return errs.NewUnauthorizedError("operation not permitted")
}
