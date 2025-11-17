package service

import (
	"errors"
	"time"

	"authorization-vendor/internal/model"
	"authorization-vendor/internal/repository"
)

type AuthorizationVendorService interface {
	AutorizarVendedor(idVendedor int, usuarioRRHH string) (bool, error)
	RevocarVendedor(idVendedor int, usuarioRRHH string) (bool, error)
	ConsultarEstado(idVendedor int) (bool, error)
	Historial(idVendedor int) ([]model.HistorialAutorizacion, error)
	ListarVendedores() ([]model.Vendedor, error)
}

type authorizationVendorService struct {
	repo repository.AuthorizationVendorRepository
}

func NewAuthorizationVendorService(repo repository.AuthorizationVendorRepository) AuthorizationVendorService {
	return &authorizationVendorService{repo: repo}
}

func (s *authorizationVendorService) AutorizarVendedor(idVendedor int, usuarioRRHH string) (bool, error) {
	return s.changeState(idVendedor, true, usuarioRRHH)
}

func (s *authorizationVendorService) RevocarVendedor(idVendedor int, usuarioRRHH string) (bool, error) {
	return s.changeState(idVendedor, false, usuarioRRHH)
}

func (s *authorizationVendorService) ConsultarEstado(idVendedor int) (bool, error) {
	estado, ok := s.repo.ObtenerEstado(idVendedor)
	if !ok {
		return false, errors.New("vendedor no encontrado")
	}
	return estado, nil
}

func (s *authorizationVendorService) Historial(idVendedor int) ([]model.HistorialAutorizacion, error) {
	return s.repo.ListarHistorial(idVendedor)
}

func (s *authorizationVendorService) ListarVendedores() ([]model.Vendedor, error) {
	return s.repo.ListarVendedores()
}

func (s *authorizationVendorService) changeState(idVendedor int, nuevoEstado bool, usuario string) (bool, error) {
	if idVendedor <= 0 {
		return false, errors.New("idVendedor invalido")
	}
	actual, ok := s.repo.ObtenerEstado(idVendedor)
	if !ok {
		actual = false
	}
	if err := s.repo.CambiarEstado(idVendedor, nuevoEstado); err != nil {
		return false, err
	}
	h := model.HistorialAutorizacion{
		IDHistorial:    int(time.Now().UnixNano()),
		IDVendedor:     idVendedor,
		EstadoAnterior: actual,
		EstadoNuevo:    nuevoEstado,
		FechaCambio:    time.Now(),
		UsuarioRRHH:    usuario,
	}
	if err := s.repo.GuardarHistorial(h); err != nil {
		return false, err
	}
	return nuevoEstado, nil
}
