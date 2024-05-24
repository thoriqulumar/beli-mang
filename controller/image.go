package controller

import (
	"beli-mang/model"
	"beli-mang/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type ImageController struct {
	svc service.ImageService
}

func NewImageController(svc service.ImageService) *ImageController {
	return &ImageController{
		svc: svc,
	}
}

func (ctr *ImageController) PostImage(c echo.Context) error {
	// Read form data including uploaded file
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "params not valid"})
	}

	file := form.File["file"][0]

	fileName := file.Filename
	if !strings.HasSuffix(strings.ToLower(fileName), ".jpg") && !strings.HasSuffix(strings.ToLower(fileName), ".jpeg") {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Only JPG/JPEG files are allowed"})
	}

	fileSize := file.Size
	if fileSize > 2*1024*1024 || fileSize < 10*1024 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "File size must be between 10KB and 2MB"})
	}

	urlChan := ctr.svc.UploadImage(file)

	url := <-urlChan
	if url == "" {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal server error"})
	}

	jsonResponse := fmt.Sprintf(`{
		"message": "File uploaded successfully",
		"data": {
			"imageUrl": "%s"
		}
	}`, url)

	var resp model.PostImageResponse
	err = json.Unmarshal([]byte(jsonResponse), &resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, resp)
}
