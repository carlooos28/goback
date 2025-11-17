package model

type Persona struct {
	IDPersona     int    `json:"idPersona"`
	TipoDocumento string `json:"tipoDocumento"`
	NumeroDoc     string `json:"numeroDocumento"`
	Nombres       string `json:"nombres"`
	Apellidos     string `json:"apellidos"`
	Direccion     string `json:"direccion"`
	Correo        string `json:"correo"`
	Celular       string `json:"celular"`
}
