package model

import "time"

type HistorialAutorizacion struct {
	IDHistorial    int       `json:"idHistorial"`
	IDVendedor     int       `json:"idVendedor"`
	EstadoAnterior bool      `json:"estadoAnterior"`
	EstadoNuevo    bool      `json:"estadoNuevo"`
	FechaCambio    time.Time `json:"fechaCambio"`
	UsuarioRRHH    string    `json:"usuarioRRHH"`
}
