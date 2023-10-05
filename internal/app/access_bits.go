package app

type accesBits [4]byte
type AccessBits accesBits

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
