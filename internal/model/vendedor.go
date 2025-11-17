package model

import "time"

type Vendedor struct {
	IDVendedor       int       `json:"idVendedor"`
	IDPersona        int       `json:"idPersona"`
	Sueldo           float64   `json:"sueldo"`
	EstadoAutorizado bool      `json:"estadoAutorizado"`
	FechaAutorizado  time.Time `json:"fechaAutorizacion,omitempty"`
	FechaRevocado    time.Time `json:"fechaRevocacion,omitempty"`
}
