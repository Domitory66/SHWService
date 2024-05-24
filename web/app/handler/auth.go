package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) mainPage(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

func (h *Handler) signIN(c *gin.Context) {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*5)
	//defer cancel()
	//resp, err := h.Services.AuthClient.UserEnter(ctx, &auth.EnterRequest{c.Param("email"), c.Param("password")})
	c.Redirect(301, "/api/listCameras")
}

// TODO Перенаправление на сайт для получения аутентификатора от YandexAPI
func (h *Handler) signUp(c *gin.Context) {
	c.Redirect(301, "https://login.yandex.ru/info?&format=jwt&jwt_secret=")
}

func (h *Handler) logout(c *gin.Context) {
	c.Redirect(301, "/")
}
