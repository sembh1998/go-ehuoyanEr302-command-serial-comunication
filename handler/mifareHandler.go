package handler

import (
	"github.com/sembh1998/go-ehuoyanEr302-command-serial-comunication/commands"
)

type MifareReader struct{}

func NewMifareReader() commands.Commands {
	return &MifareReader{}
}

func (m *MifareReader) Read(auth commands.Auth) (commands.Block, error) {
	return commands.Block{}, nil
}

func (m *MifareReader) Write(auth commands.Auth, block commands.Block) error {
	return nil
}

func (m *MifareReader) ChangeKeys(auth commands.Auth, key_a, key_b commands.Key, accestype commands.AccessBits) error {
	return nil
}
