package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"rest-api-go/internal/apperrors"
	userEntity "rest-api-go/internal/entities/user"
	"rest-api-go/internal/handlers/interfaces"
	"rest-api-go/internal/service"

	"rest-api-go/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

const (
	usersUrl = "/users"
	userUrl  = "/users/:uuid"
)

type UserHandler struct {
	logger      *logging.Logger
	userService service.UserService
}

func NewUserHandler(logger *logging.Logger, userService service.UserService) interfaces.Handler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
	}
}

func (h *UserHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersUrl, apperrors.Middleware(h.GetAll))
	router.HandlerFunc(http.MethodGet, userUrl, apperrors.Middleware(h.GetUserByUUID))
	router.HandlerFunc(http.MethodPost, usersUrl, apperrors.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodPut, userUrl, apperrors.Middleware(h.UpdateUser))
	router.HandlerFunc(http.MethodDelete, userUrl, apperrors.Middleware(h.DeleteUser))

}
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	users, err := h.userService.FindAll(r.Context())
	if err != nil {
		return err
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to marshal users", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(usersJSON)
	return nil
}
func (h *UserHandler) GetUserByUUID(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("GET USER")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Debug("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userUUID := params.ByName("uuid")

	user, err := h.userService.FindOne(r.Context(), userUUID)
	if err != nil {
		return err
	}

	h.logger.Debug("marshal user")
	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshall user. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("CREATE USER")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Debug("decode create user dto")
	var crUser userEntity.CreateUserDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&crUser); err != nil {
		return apperrors.BadRequestError("invalid JSON scheme. check swagger API")
	}

	userUUID, err := h.userService.Create(r.Context(), crUser)
	if err != nil {
		return err
	}
	w.Header().Set("Location", fmt.Sprintf("%s/%s", usersUrl, userUUID))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(userUUID))

	return nil
}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("PARTIALLY UPDATE USER")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userUUID := params.ByName("uuid")

	h.logger.Debug("decode update user dto")
	var updUser userEntity.UpdateUserDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&updUser); err != nil {
		return apperrors.BadRequestError("invalid JSON scheme. check swagger API")
	}
	updUser.ID = userUUID

	err := h.userService.Update(r.Context(), updUser)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("DELETE USER")
	w.Header().Set("Content-Type", "application/json")

	h.logger.Debug("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userUUID := params.ByName("uuid")

	err := h.userService.Delete(r.Context(), userUUID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(userUUID))

	return nil
}
