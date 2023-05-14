package handler

import (
	"syscall"

	"github.com/sembh1998/go-ehuoyanEr302-command-serial-comunication/commands"
)

type MifareReader struct {
	handle syscall.Handle
	keyA   commands.Key
	keyB   commands.Key
}

func NewMifareReader(key_a, key_b commands.Key, handler syscall.Handle) commands.Commands {
	return &MifareReader{
		keyA:   key_a,
		keyB:   key_b,
		handle: handler,
	}
}

func (m *MifareReader) Read() ([]byte, error) {
	return nil, nil
}

func (m *MifareReader) Write(data []byte) error {
	return nil
}

func (m *MifareReader) ChangeKeys(old_key_a, old_key_b commands.Key) error {
	return nil
}
