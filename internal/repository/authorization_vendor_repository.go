package repository

import "authorization-vendor/internal/model"

type AuthorizationVendorRepository interface {
	ObtenerEstado(idVendedor int) (bool, bool)
	CambiarEstado(idVendedor int, estado bool) error
	GuardarHistorial(h model.HistorialAutorizacion) error
	ListarHistorial(idVendedor int) ([]model.HistorialAutorizacion, error)
	ListarVendedores() ([]model.Vendedor, error)
}
