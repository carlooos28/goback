package controller

import (
	"net/http"

	"authorization-vendor/internal/service"
)

type AuthorizationVendorController struct {
	service service.AuthorizationVendorService
}

func NewAuthorizationVendorController(s service.AuthorizationVendorService) *AuthorizationVendorController {
	return &AuthorizationVendorController{service: s}
}

// Expose service to router handlers; thin layer kept for clarity with the class diagram.
func (c *AuthorizationVendorController) Service() service.AuthorizationVendorService {
	return c.service
}

func (c *AuthorizationVendorController) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
