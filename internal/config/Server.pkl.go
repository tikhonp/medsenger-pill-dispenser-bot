// Code generated from Pkl module `config`. DO NOT EDIT.
package config

type Server struct {
	// The port to listen on.
	Port uint16 `pkl:"port"`

	// Sets server to debug mode.
	Debug bool `pkl:"debug"`

	// Medsenger Agent secret key.
	MedsengerAgentKey string `pkl:"medsengerAgentKey"`
}
