package models

// Config - config structure
type Config struct {
	UDPPath    string `json:"udp_path"`
	FPS        int    `json:"fps"`
	MaxBufSize int    `json:"max_buf_size"`
}

// Package - contains parts data
type Package []PartPackage

/*
PartPackage - contains the data of the part of the package
	- PartPackage - contains:

		- ID <string> - part's ID of package.

		- Part <int> - number of part.

		- Data <[]byte> - data of package's part.

		- End <bool> - determines if the last part is.
*/
type PartPackage struct {
	ID   string `json:"id"`
	Part int    `json:"part"`
	Data []byte `json:"data"`
	End  bool   `json:"end"`
}
