package api

import (
	"github.com/gin-gonic/gin"
	"practice/model"
	"practice/service/shorten"
)

func Shorten(c *gin.Context) {
	var req model.GetShortLinkReq
	err := c.ShouldBind(&req)
	if err != nil {
		ErrorResponse(c, 500, err.Error(), nil)
		return
	}
	short, err := shorten.GetShortLink(req.URL)
	if err != nil {
		ErrorResponse(c, 500, err.Error(), nil)
		return
	}
	SuccessResponse(c, 200, "success", short)
}

func Redirect(c *gin.Context) {
	shortLink := c.Param("short_link")
	longURL, err := shorten.GetLongURL(shortLink)
	if err != nil {
		ErrorResponse(c, 500, err.Error(), nil)
		return
	}
	c.Redirect(302, longURL)
}
