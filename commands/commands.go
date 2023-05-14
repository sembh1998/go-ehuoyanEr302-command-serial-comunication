package commands

type Key [12]byte

type Commands interface {
	Read() ([]byte, error)
	Write(data []byte) error
	ChangeKeys(old_key_a, old_key_b Key) error
}
