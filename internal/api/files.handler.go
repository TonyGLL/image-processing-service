package api

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	db "github.com/TonyGLL/image-processing-service/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

func (s *Server) uploadImageHandler(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File not received"})
		return
	}
	defer file.Close()

	filename := header.Filename
	ext := strings.ToLower(filepath.Ext(filename))
	fmt.Println(ext)

	// Allowed file extensions
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if !allowedExtensions[ext] {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file format"})
		return
	}

	// Create images directory if it doesn't exist
	saveDir := "./images"
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	// Save the uploaded file
	savePath := filepath.Join(saveDir, filename)
	out, err := os.Create(savePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the image"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write the image"})
		return
	}

	// Open the saved file to read metadata
	savedFile, err := os.Open(savePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open saved image"})
		return
	}
	defer savedFile.Close()

	// Detect MIME type
	buffer := make([]byte, 512)
	_, err = savedFile.Read(buffer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image data"})
		return
	}
	contentType := http.DetectContentType(buffer)

	// Get file info
	fileInfo, err := savedFile.Stat()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
		return
	}

	// Build image URL (adjust the domain/protocol as needed)
	imageURL := fmt.Sprintf("%s://%s/images/%s", ctx.Request.URL.Scheme, ctx.Request.Host, filename)
	if imageURL[:5] != "http:" && imageURL[:6] != "https:" {
		imageURL = fmt.Sprintf("http://%s/images/%s", ctx.Request.Host, filename)
	}

	imageId, err := s.store.CreateImage(ctx, imageURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
		return
	}

	params := db.CreateImageOptionsParams{
		ResizeWidth:  0,
		ResizeHeight: 0,
		CropWidth:    0,
		CropHeight:   0,
		CropX:        0,
		CropY:        0,
		Rotate:       0,
		Format:       "a",
		Grayscale:    false,
		Sepia:        false,
		ImageID:      imageId,
	}

	if err := s.store.CreateImageOptions(ctx, params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Image uploaded successfully",
		"url":     imageURL,
		"metadata": gin.H{
			"name":         filename,
			"size_bytes":   fileInfo.Size(),
			"content_type": contentType,
		},
	})
}

func (s *Server) listImagesHandler(ctx *gin.Context) {
	params := db.GetAllImagesParams{
		Offset: 0,
		Limit:  10,
	}
	data, err := s.store.GetAllImages(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"images": data})
}

type getImageByIdReq struct {
	ID int32 `uri:"id" binding:"required"`
}

func (s *Server) getImageByIdHandler(ctx *gin.Context) {
	var req getImageByIdReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	data, err := s.store.GetImageById(ctx, int32(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"images": data})
}
