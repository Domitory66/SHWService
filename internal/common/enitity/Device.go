package entity

type Device struct {
	Name           string
	Type           string
	TypeConnection string
	IsConnected    bool
	Params         interface{}
}
