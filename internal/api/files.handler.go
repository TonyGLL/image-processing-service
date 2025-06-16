package api

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"image"
	"log"
	"path/filepath"
	"strings"

	db "github.com/TonyGLL/image-processing-service/internal/db/sqlc"
	"github.com/disintegration/imaging"
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

type TransformationPayload struct {
	Resize *struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"resize,omitempty"`
	Crop *struct {
		Width  int `json:"width"`
		Height int `json:"height"`
		X      int `json:"x"`
		Y      int `json:"y"`
	} `json:"crop,omitempty"`
	Rotate  *int    `json:"rotate,omitempty"`   // Pointer to allow optional
	Format  *string `json:"format,omitempty"`   // Pointer to allow optional
	Filters *struct {
		Grayscale *bool `json:"grayscale,omitempty"` // Pointer to allow optional
		Sepia     *bool `json:"sepia,omitempty"`     // Pointer to allow optional
	} `json:"filters,omitempty"`
}

type ImageTransformRequest struct {
	Transformations TransformationPayload `json:"transformations" binding:"required"`
}

type getImageIdURI struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *Server) transformImageHandler(ctx *gin.Context) {
	var uri getImageIdURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req ImageTransformRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	imageData, err := s.store.GetImageById(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("image with ID %d not found", uri.ID)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// Derive local file path from URL. This is a simplification and might need adjustment.
	// Assumes URL is like "http://localhost:3000/images/filename.jpg"
	// and files are stored in "./images/"
	filename := filepath.Base(imageData.Url)
	localImagePath := filepath.Join("./images", filename)

	if req.Transformations.Resize != nil && (req.Transformations.Resize.Width > 0 || req.Transformations.Resize.Height > 0) {
		log.Printf("Resizing image ID %d to width %d, height %d", uri.ID, req.Transformations.Resize.Width, req.Transformations.Resize.Height)

		srcImage, err := imaging.Open(localImagePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to open image %s: %v", localImagePath, err)))
			return
		}

		dstImage := imaging.Resize(srcImage, req.Transformations.Resize.Width, req.Transformations.Resize.Height, imaging.Lanczos)

		// Overwrite the original image
		err = imaging.Save(dstImage, localImagePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to save resized image %s: %v", localImagePath, err)))
			return
		}
		log.Printf("Successfully resized and saved image %s", localImagePath)

		// Update database
		params := db.UpdateImageResizeOptionsParams{
			ResizeWidth:  int32(req.Transformations.Resize.Width),
			ResizeHeight: int32(req.Transformations.Resize.Height),
			ImageID:      uri.ID,
		}
		err = s.store.UpdateImageResizeOptions(ctx, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update image resize options in DB: %v", err)))
			return
		}
		log.Printf("Successfully updated resize options for image ID %d in DB", uri.ID)
	}

	if req.Transformations.Crop != nil && req.Transformations.Crop.Width > 0 && req.Transformations.Crop.Height > 0 {
		log.Printf("Cropping image ID %d to width %d, height %d at x:%d, y:%d",
			uri.ID, req.Transformations.Crop.Width, req.Transformations.Crop.Height, req.Transformations.Crop.X, req.Transformations.Crop.Y)

		currentImage, err := imaging.Open(localImagePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to open image %s for cropping: %v", localImagePath, err)))
			return
		}

		cropRect := image.Rect(
			req.Transformations.Crop.X,
			req.Transformations.Crop.Y,
			req.Transformations.Crop.X+req.Transformations.Crop.Width,
			req.Transformations.Crop.Y+req.Transformations.Crop.Height,
		)
		dstImage := imaging.Crop(currentImage, cropRect)

		err = imaging.Save(dstImage, localImagePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to save cropped image %s: %v", localImagePath, err)))
			return
		}
		log.Printf("Successfully cropped and saved image %s", localImagePath)

		params := db.UpdateImageCropOptionsParams{
			CropWidth:  int32(req.Transformations.Crop.Width),
			CropHeight: int32(req.Transformations.Crop.Height),
			CropX:      int32(req.Transformations.Crop.X),
			CropY:      int32(req.Transformations.Crop.Y),
			ImageID:    uri.ID,
		}
		err = s.store.UpdateImageCropOptions(ctx, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update image crop options in DB: %v", err)))
			return
		}
		log.Printf("Successfully updated crop options for image ID %d in DB", uri.ID)
	}

	// After all transformations, fetch the updated image data
	updatedImageData, err := s.store.GetImageById(ctx, uri.ID) // Refetch to get updated options
	if err != nil {
		// Handle error, though unlikely if previous steps succeeded
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message":           "Image transformations applied successfully.",
		"image_id":          uri.ID,
		"updated_details": updatedImageData,
	})
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
