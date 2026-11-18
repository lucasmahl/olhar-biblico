package models

type AuthenticationData struct {
	Authorized    bool   `json:"authorized,omitempty"`
	Exp           int64  `json:"exp,omitempty"`
	UsuarioID     string `json:"usuario_id,omitempty"`
	RemoteAddress string `json:"remote_address,omitempty"`
}
