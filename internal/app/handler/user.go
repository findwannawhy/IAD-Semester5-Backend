package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/findwannawhy/IAD-Semester5/internal/app/ds"
	"github.com/findwannawhy/IAD-Semester5/internal/app/dto"
	"github.com/findwannawhy/IAD-Semester5/internal/app/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateUser(ctx *gin.Context) {
	var userJSON ds.User
	if err := ctx.BindJSON(&userJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Валидация обязательных полей
	if userJSON.Login == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("field 'login' is required"))
		return
	}
	if userJSON.Password == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("field 'password' is required"))
		return
	}

	user, err := h.Repository.CreateUser(userJSON)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			h.errorHandler(ctx, http.StatusConflict, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Header("Location", fmt.Sprintf("/user/%v", user.ID))
	ctx.JSON(http.StatusCreated, dto.UserResponse{
		ID: user.ID,
		Login: user.Login,
		IsModerator: user.IsModerator,
	})
}

func (h *Handler) SignIn(ctx *gin.Context) {
	var userJSON ds.User
	if err := ctx.BindJSON(&userJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Валидация обязательных полей
	if userJSON.Login == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("field 'login' is required"))
		return
	}
	if userJSON.Password == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("field 'password' is required"))
		return
	}

	user, err := h.Repository.SignIn(userJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("invalid login or password"))
		} else if err.Error() == "wrong password" {
			h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("invalid login or password"))
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, dto.UserResponse{
		ID: user.ID,
		Login: user.Login,
		IsModerator: user.IsModerator,
	})
}

func (h *Handler) GetProfile(ctx *gin.Context) {
	userID := h.Repository.GetUserID()
	if userID == 0 {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("user not authenticated"))
		return
	}

	user, err := h.Repository.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.JSON(http.StatusOK, dto.UserResponse{
		ID: user.ID,
		Login: user.Login,
		IsModerator: user.IsModerator,
	})
}

func (h *Handler) UpdateProfile(ctx *gin.Context) {
	userID := h.Repository.GetUserID()
	if userID == 0 {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("user not authenticated"))
		return
	}

	var userJSON ds.User
	if err := ctx.BindJSON(&userJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.UpdateProfile(userID, userJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.JSON(http.StatusOK, dto.UserResponse{
		ID: user.ID,
		Login: user.Login,
		IsModerator: user.IsModerator,
	})
}

func (h *Handler) SignOut(ctx *gin.Context) {
	//h.Repository.SignOut()
	ctx.JSON(http.StatusOK, gin.H{
		"status": "signed_out",
	})
}