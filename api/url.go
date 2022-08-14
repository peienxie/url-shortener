package api

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/peienxie/url-shortener/shorten"
)

var ErrBadLongURL = errors.New("'long_url' parse error")
var ErrURLSchemeNotSupport = errors.New("scheme of 'long_url' only supports http or https")
var ErrURLHostEmpty = errors.New("host of 'long_url' is empty")
var ErrExpireTimePast = errors.New("'expire_time' is past")
var ErrShortURLNotFound = errors.New("'short_url' is past")

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

type CreateShortURLReq struct {
	LongURL    string    `json:"long_url" binding:"required"`
	ExpireTime time.Time `json:"expire_at" binding:"required" time_format:"2006-01-02T15:04:05Z"`
}

type CreateShortURLRsp struct {
	ShortURL   string    `json:"short_url"`
	ExpireTime time.Time `json:"expire_at"`
}

func (s *Server) createShortURL(c *gin.Context) {
	var err error
	var req CreateShortURLReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check long_url parameter is valid url
	u, err := url.Parse(req.LongURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(ErrBadLongURL))
		return
	} else if u.Scheme != "http" && u.Scheme != "https" {
		c.JSON(http.StatusBadRequest, errorResponse(ErrURLSchemeNotSupport))
		return
	} else if u.Host == "" {
		c.JSON(http.StatusBadRequest, errorResponse(ErrURLHostEmpty))
		return
	}

	now := time.Now()
	if now.After(req.ExpireTime) {
		c.JSON(http.StatusBadRequest, errorResponse(ErrExpireTimePast))
		return
	}

	shortURL := shorten.ShortenByHash(req.LongURL)
	err = s.store.SaveURL(c, shortURL, req.LongURL, req.ExpireTime.Sub(now))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, CreateShortURLRsp{
		ShortURL:   shortURL,
		ExpireTime: req.ExpireTime,
	})
}

type RedirectShortURLReq struct {
	ShortURL string `uri:"short_url" binding:"required"`
}

func (s *Server) redirectShortURL(c *gin.Context) {
	var req RedirectShortURLReq
	var err error
	err = c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	longURL, err := s.store.LoadURL(c, req.ShortURL)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.Redirect(http.StatusFound, longURL)
}
