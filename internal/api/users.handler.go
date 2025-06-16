package api

import (
	"database/sql"
	"net/http"

	db "github.com/TonyGLL/image-processing-service/internal/db/sqlc"
	"github.com/TonyGLL/image-processing-service/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

/* CREATE USERS */
type CreateUserRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerUserResponse struct {
	UserID   int32  `json:"user_id"`
	UserName string `json:"username"`
	Token    string `json:"token"`
}

// Create user	godoc
// @Summary Create user
// @Description Create user by ID
// @Tags Users
// @Accept json
// @Produce application/json
// @Param			key	path		createUserRequest		true	"Site KEY"
// @in header
// @name Authorization
// @Success 200 {object} string
// @Failure		400			{string}	gin.H	"StatusBadRequest"
// @Failure		404			{string}	gin.H	"StatusNotFound"
// @Failure		500			{string}	gin.H	"StatusInternalServerError"
// @Router /users [post]
func (s *Server) createUserHandler(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := s.store.CreateUser(ctx, req.UserName)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	createPasswordParams := db.CreatePasswordParams{
		Value:  string(hashedPassword),
		UserID: userId,
	}

	err = s.store.CreatePassword(ctx, createPasswordParams)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	token, err := utils.GenerateToken(req.UserName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := registerUserResponse{
		UserID:   userId,
		UserName: req.UserName,
		Token:    token,
	}
	ctx.JSON(http.StatusOK, response)
}
