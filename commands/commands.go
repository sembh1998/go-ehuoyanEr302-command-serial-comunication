package commands

type Key [6]byte
type Block [16]byte
type Sector [4]Block
type accesBits [4]byte
type BlockAddress uintptr
type AccessBits accesBits

type Auth struct {
	Key     Key
	KeyType keyType
	Address BlockAddress
}

type keyType byte

const (
	KeyAType keyType = 0x60
	KeyBType keyType = 0x61
)

var (
	// 787788
	// KeyA Only Read, KeyB Read and Write
	ABits_keyaOR_keybRW accesBits = accesBits{0x78, 0x77, 0x88, 0xFF}

	// ff0780
	// KeyA Read and Write, KeyB Disabled
	ABits_keyaRW_keybDIS accesBits = accesBits{0xFF, 0x07, 0x80, 0xFF}

	// 0f00ff
	// KeyA Disabled, KeyB Read and Write
	ABits_keyaDIS_keybRW accesBits = accesBits{0x0F, 0x00, 0xFF, 0xFF}
)

type Commands interface {
	Read(auth Auth) (Block, error)
	Write(auth Auth, block Block) error
	ChangeKeys(auth Auth, key_a, key_b Key, accestype AccessBits) error
}
