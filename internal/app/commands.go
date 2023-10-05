package app

type Commands interface {
	Read(auth Auth) (Block, error)
	Write(auth Auth, block Block) error
	ChangeKeys(auth Auth, key_a, key_b Key, accestype AccessBits) error
}
