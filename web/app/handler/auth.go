package handler

import (
	"SmartHomeWebCam/SHWService/web/api/auth"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) mainPage(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

func (h *Handler) signIN(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := h.Services.AuthClient.UserEnter(ctx, &auth.EnterRequest{Email: c.Param("email"), Pass: c.Param("password")})
	if err != nil {
		switch status.Code(err) {
		default:
			c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		}
		return
	}

	c.SetSameSite(http.SameSiteDefaultMode)
	c.SetCookie("Auth", resp.GetToken(), 3600*24*30, "", "", false, true)
	c.Redirect(301, "/api/listCameras")
}

func (h *Handler) signUp(c *gin.Context) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resp, err := h.Services.AuthClient.Registration(ctx, &emptypb.Empty{})
	if err != nil {
		switch status.Code(err) {
		default:
			c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		}
		return
	}
	c.SetSameSite(http.SameSiteDefaultMode)
	c.SetCookie("Auth", resp.GetToken(), 3600*24*30, "", "", false, true)
	c.Redirect(301, "/api/listCameras")
}

func (h *Handler) logout(c *gin.Context) {
	c.Set("Auth", "")
	c.Redirect(301, "/")
}
