package handler

import (
	gen "SmartHomeWebCam/SHWService/web/api/camera"
	"SmartHomeWebCam/SHWService/web/api/video"
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Camera struct {
	Ip   string
	Port string
	Name string
}

func (h *Handler) getListCameras(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*5)
	defer cancel()
	resp, err := h.Services.CameraWorkerClient.GetAllCameras(ctx, &gen.GetAllCamerasRequest{UserID: 0})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.HTML(200, "listCameras.html", nil)
		case codes.Unavailable:
		default:
			c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		}
		return
	}
	var cams []*Camera
	for i := range resp.Cameras {
		cams = append(cams, &Camera{Ip: resp.Cameras[i].Ip, Port: resp.Cameras[i].Port, Name: resp.Cameras[i].Name})
	}
	log.Print(resp.Cameras)
	c.HTML(200, "listCameras.html", gin.H{"cameras": cams})
}

func (h *Handler) getCameraView(c *gin.Context) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// _, err := h.Services.CameraWorkerClient.GetCameraByPortAndIp(ctx, &gen.GetCameraRequest{Ip: c.Param("ip"), Port: c.Param("port")})
	// if err != nil {
	// 	switch status.Code(err) {
	// 	case codes.NotFound:
	// 		c.Redirect(301, "/api/listCameras/notFound")
	// 		return
	// 	}
	// 	log.Fatalf(err.Error())
	// }
	c.HTML(200, "camera.html", gin.H{"ip": c.Param("ip"), "port": c.Param("port")})
}

func (h *Handler) showFormAdd(c *gin.Context) {
	c.HTML(200, "addCamera.html", nil)
}

func (h *Handler) showFormNotFound(c *gin.Context) {
	c.HTML(200, "notFound.html", nil)
}

func (h *Handler) Video(c *gin.Context) {
	ctx, cancel := context.WithCancel(c)
	c.Writer.Header().Set("Content-type", "multipart/x-mixed-replace; boundary=frame")
	mutex := &sync.Mutex{}
	data := ""
	defer cancel()
	for {
		mutex.Lock()
		resp, err := h.Services.VideoStreamClient.GetVideoFromCamera(ctx, &video.ImageRequest{Ip: c.Param("ip"), Port: c.Param("port")})
		if err != nil {
			switch status.Code(err) {
			case codes.Aborted:
				return
			}
			return
		}
		frame := string(resp.Image)
		data = "--frame\r\n Content-Type: image/jpeg\r\n\r\n" + frame + "\r\n\r\n"

		log.Print(len(data))
		mutex.Unlock()
		// time.Sleep(10 * time.Millisecond)
		_, err = c.Writer.Write([]byte(data))
		if err != nil {
			ctx.Done()
			h.Services.VideoStreamClient.StopVideoStream(ctx, &video.StopRequest{Ip: c.Param("ip"), Port: c.Param("port")})
			return
		}
	}
}

func (h *Handler) SetProcess(c *gin.Context) {
	//TODO флаг включения режима распознавания жестов для данной камеры
}

func (h *Handler) addCamera(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	respAdd, err := h.Services.CameraWorkerClient.AddCamera(ctx, &gen.AddCameraRequest{
		UserID: 0,
		Camera: &gen.Camera{
			Name:     c.Request.FormValue("name"),
			Ip:       c.Request.FormValue("ip"),
			Port:     c.Request.FormValue("port"),
			Protocol: c.Request.FormValue("protocol"),
			Filename: c.Request.FormValue("file"),
		},
	})
	if err != nil {
		switch status.Code(err) {
		case codes.AlreadyExists:
			c.Redirect(301, "/api/listCameras/")
		default:
			c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		}
		return
	}
	if !respAdd.Saved {
		c.JSON(http.StatusBadGateway, gin.H{"message": "Camera not saved"})
		return
	}

	c.Redirect(301, "/api/listCameras/")
}

func (h *Handler) deleteCamera(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Second*2)
	defer cancel()

	resp, err := h.Services.CameraWorkerClient.DeleteCamera(ctx, &gen.DeleteCameraRequest{UserID: 0, Camera: &gen.Camera{Ip: c.Param("ip"), Port: c.Param("port")}})
	if err != nil {
		log.Fatal(err)
	}
	if !resp.Deleted {
		log.Fatal("Not deleted")
	}
	c.Redirect(301, "/api/listCameras/")
}
