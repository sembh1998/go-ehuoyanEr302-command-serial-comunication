package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/google/uuid"
	"syreclabs.com/go/faker"
)

const (
	libPath             = "MasterRD.dll"          // replace with the name of your DLL
	libVerFuncName      = "lib_ver"               // name of the lib_ver function
	rfInitComFuncName   = "rf_init_com"           // name of the rf_init_com function
	rfClosePortFuncName = "rf_ClosePort"          // name of the rf_ClosePort function
	rfLightFuncName     = "rf_light"              // name of the rf_light function
	rfBeepFuncName      = "rf_beep"               // name of the rf_beep function
	rfRequestFuncName   = "rf_request"            // name of the rf_request function
	rfAnticollFuncName  = "rf_anticoll"           // name of the rf_anticoll function
	rfSelectFuncName    = "rf_select"             // name of the rf_select function
	rfM1AuthFuncName    = "rf_M1_authentication2" // name of the rf_M1_authentication2 function
	rfM1ReadFuncName    = "rf_M1_read"            // name of the rf_M1_read function
	rfM1WriteFuncName   = "rf_M1_write"           // name of the rf_M1_write function
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

var rfClosePort uintptr

func handleFailedCallWithMessage(ret uintptr, msg string) {
	if ret != 0 {
		ret2, _, _ := syscall.SyscallN(rfClosePort)
		log.Printf("rf_ClosePort returned %d", ret2)
		panic(fmt.Sprintf("%s failed with error code %d", msg, ret))
	}
}

func main() {
	// Load the DLL
	handle, err := syscall.LoadLibrary(libPath)
	handleError(err)
	defer syscall.FreeLibrary(handle)

	// Get the lib_ver function
	libVer, err := syscall.GetProcAddress(handle, libVerFuncName)
	handleError(err)

	// Get the rf_init_com function
	rfInitCom, err := syscall.GetProcAddress(handle, rfInitComFuncName)
	handleError(err)

	// Get the rf_ClosePort function
	rfClosePort, err = syscall.GetProcAddress(handle, rfClosePortFuncName)
	handleError(err)

	// Get the rf_light function
	rfLight, err := syscall.GetProcAddress(handle, rfLightFuncName)
	handleError(err)

	// Get the rf_beep function
	rfBeep, err := syscall.GetProcAddress(handle, rfBeepFuncName)
	handleError(err)

	// Get the rf_request function
	rfRequest, err := syscall.GetProcAddress(handle, rfRequestFuncName)
	handleError(err)

	// Get the rf_anticoll function
	rfAnticoll, err := syscall.GetProcAddress(handle, rfAnticollFuncName)
	handleError(err)

	// Get the rf_select function
	rfSelect, err := syscall.GetProcAddress(handle, rfSelectFuncName)
	handleError(err)

	// Get the rf_M1_authentication2 function
	rfM1Auth, err := syscall.GetProcAddress(handle, rfM1AuthFuncName)
	handleError(err)

	// Get the rf_M1_read function
	rfM1Read, err := syscall.GetProcAddress(handle, rfM1ReadFuncName)
	handleError(err)

	// Get the rf_M1_write function
	rfM1Write, err := syscall.GetProcAddress(handle, rfM1WriteFuncName)
	handleError(err)

	// Call the lib_ver function
	var ver uint32
	ret, _, _ := syscall.SyscallN(uintptr(libVer), uintptr(unsafe.Pointer(&ver)))
	handleFailedCallWithMessage(ret, libVerFuncName)
	fmt.Printf("DLL version: %d\n", ver)

	// Call the rf_init_com function
	ret, _, _ = syscall.SyscallN(uintptr(rfInitCom), 3, 115200)
	handleFailedCallWithMessage(ret, rfInitComFuncName)
	fmt.Println("Connected successfully")

	// Call the rf_light function
	ret, _, _ = syscall.SyscallN(uintptr(rfLight), 0, 0)
	handleFailedCallWithMessage(ret, rfLightFuncName)
	fmt.Println("LED turned off")
	// Wait for 1 second
	time.Sleep(300 * time.Millisecond)

	// Call the rf_light function
	ret, _, _ = syscall.SyscallN(uintptr(rfLight), 0, 2)
	handleFailedCallWithMessage(ret, rfLightFuncName)
	fmt.Println("LED turned red")
	// Wait for 1 second
	time.Sleep(300 * time.Millisecond)

	// Call the rf_light function
	ret, _, _ = syscall.SyscallN(uintptr(rfLight), 0, 1)
	handleFailedCallWithMessage(ret, rfLightFuncName)
	fmt.Println("LED turned blue")
	// Wait for 1 second
	time.Sleep(300 * time.Millisecond)

	// Call the rf_beep function
	ret, _, _ = syscall.SyscallN(uintptr(rfBeep), 0, 10)
	handleFailedCallWithMessage(ret, rfBeepFuncName)
	fmt.Println("Beeped")
	time.Sleep(300 * time.Millisecond)

	// Call the rf_beep function
	ret, _, _ = syscall.SyscallN(uintptr(rfBeep), 0, 10)
	handleFailedCallWithMessage(ret, rfBeepFuncName)
	fmt.Println("Beeped")
	time.Sleep(300 * time.Millisecond)

	// Call the rf_request function
	var mode byte = 0x52
	var tagType uint32
	ret, _, _ = syscall.SyscallN(uintptr(rfRequest), 0, uintptr(unsafe.Pointer(&mode)), uintptr(unsafe.Pointer(&tagType)))
	handleFailedCallWithMessage(ret, rfRequestFuncName)
	fmt.Printf("Tag type: %d\n", tagType)

	// Call the rf_anticoll function
	var snrLen uint32
	var snr [4]byte
	ret, _, _ = syscall.SyscallN(uintptr(rfAnticoll), 0, uintptr(tagType), uintptr(unsafe.Pointer(&snr)), uintptr(unsafe.Pointer(&snrLen)))
	handleFailedCallWithMessage(ret, rfAnticollFuncName)
	fmt.Printf("Tag serial number length: %d\n", snrLen)
	fmt.Printf("Tag serial number: % X\n", snr)

	// Call the rf_select function
	var size uint32
	ret, _, _ = syscall.SyscallN(uintptr(rfSelect), 0, uintptr(unsafe.Pointer(&snr)), uintptr(snrLen), uintptr(unsafe.Pointer(&size)))
	handleFailedCallWithMessage(ret, rfSelectFuncName)
	fmt.Printf("Tag size: %d\n", size)

	// Call the rf_M1_authentication2 function
	// the default_key is 6 bytes long: FFFFFFFFFFFF
	// the new key_a is 6 bytes long: 15F63DAF0A38
	// the new key_b is 6 bytes long: 17292304FFA7
	// force brute attack in Trying key: 81 1A 03 00 00 00
	// BruteForceAttack(rfM1Auth)
	var blockAddr uint32 = 8
	var key_used string = "default_key"
	key_a := [6]byte{0x15, 0xF6, 0x3D, 0xAF, 0x0A, 0x38}
	key_b := [6]byte{0x17, 0x29, 0x23, 0x04, 0xFF, 0xA7}
	default_key := [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	fmt.Printf("Sector %d\n", (blockAddr+1)/4)
	fmt.Printf("Block %d\n", blockAddr)
	key_used = "key_a"
	fmt.Printf("Authenticating with %s: % X\n", key_used, key_a)
	ret, _, _ = syscall.SyscallN(uintptr(rfM1Auth), 0, 0x60, uintptr(blockAddr), uintptr(unsafe.Pointer(&key_a)))
	if ret != 0 && ret != 22 {
		handleFailedCallWithMessage(ret, rfM1AuthFuncName)
	}
	if ret == 22 {
		fmt.Println("Authentication failed with", key_used)
		key_used = "key_b"
		fmt.Printf("Authenticating with %s: % X\n", key_used, key_b)
		ret, _, _ = syscall.SyscallN(uintptr(rfM1Auth), 0, 0x61, uintptr(blockAddr), uintptr(unsafe.Pointer(&key_b)))
		if ret != 0 && ret != 22 {
			handleFailedCallWithMessage(ret, rfM1AuthFuncName)
		}
		if ret == 22 {
			fmt.Println("Authentication failed with", key_used)
			key_used = "default_key"
			fmt.Printf("Authenticating with %s: % X\n", key_used, default_key)
			ret, _, _ = syscall.SyscallN(uintptr(rfM1Auth), 0, 0x60, uintptr(blockAddr), uintptr(unsafe.Pointer(&default_key)))
			if ret != 0 && ret != 22 {
				handleFailedCallWithMessage(ret, rfM1AuthFuncName)
			}
			if ret == 22 {
				fmt.Println("Authentication failed with all keys")
				os.Exit(1)
			}
		}
	}
	fmt.Println("Authenticated")

	// Call the rf_M1_read function
	var blockData [16]byte
	var blockLen uint32
	blockAddr = 11
	fmt.Printf("Sector %d\n", (blockAddr+1)/4)
	fmt.Printf("Block %d\n", blockAddr)
	ret, _, _ = syscall.SyscallN(uintptr(rfM1Read), 0, uintptr(blockAddr), uintptr(unsafe.Pointer(&blockData)), uintptr(unsafe.Pointer(&blockLen)))
	handleFailedCallWithMessage(ret, rfM1ReadFuncName)
	fmt.Printf("Block data readed: % X\n", blockData)
	fmt.Printf("Block length: %d\n", blockLen)
	// check if content of block 2 is an uuid
	id_read, err := uuid.FromBytes(blockData[:])
	if err != nil {
		fmt.Printf("Block %v does not contain an uuid\n", blockAddr)
	} else {
		fmt.Printf("Block %v contains an uuid: %s\n", blockAddr, id_read.String())
	}
	// print the block data as a string
	fmt.Printf("Block data as string: %s\n", blockData[:])

	var blockData2 [16]byte
	var data = []byte(faker.Commerce().ProductName())
	copy(blockData2[:], data)

	if key_used == "default_key" {
		// Call the rf_M1_write function to write the new key_a
		// replace the first 6 bytes of block 3 with the new key_a

		blockData2[0] = key_a[0]
		blockData2[1] = key_a[1]
		blockData2[2] = key_a[2]
		blockData2[3] = key_a[3]
		blockData2[4] = key_a[4]
		blockData2[5] = key_a[5]
		// put the 4 bytes of the access bits 08778F
		// this means key_a is only for read and key_b is for read and write

		// put the 4 bytes of the access bits F0FF00
		// this means key_a is for read and write and key_b is only for read
		// correction: this really means that key_a can read block 0,1,2, but cannot changes the keys neither the access bits
		// key_b can read and write block 0,1,2, can change the keys but cannot changes the access bits
		// so, warning, using this access bits means the access bits cannot be changed anymore
		blockData2[6] = 0xF0
		blockData2[7] = 0xFF
		blockData2[8] = 0x00
		blockData2[9] = 0x69

		// put the 6 bytes of the new key_b
		blockData2[10] = key_b[0]
		blockData2[11] = key_b[1]
		blockData2[12] = key_b[2]
		blockData2[13] = key_b[3]
		blockData2[14] = key_b[4]
		blockData2[15] = key_b[5]

		blockAddr = 11
		fmt.Printf("blockData2 data: % X\n", blockData2)
		fmt.Printf("Sector %d\n", (blockAddr+1)/4)
		fmt.Printf("Block %d\n", blockAddr)
		ret, _, _ = syscall.SyscallN(uintptr(rfM1Write), 0, uintptr(blockAddr), uintptr(unsafe.Pointer(&blockData2)))
		handleFailedCallWithMessage(ret, rfM1WriteFuncName)
		fmt.Println("Keys written", blockData2[:])
	}
	// Call the rf_M1_write function
	// save random data to block 1
	blockAddr = 8
	new_uuid := uuid.New()
	data_to_write, _ := new_uuid.MarshalBinary()
	copy(blockData2[:], data_to_write)
	fmt.Printf("New uuid: %s\n", new_uuid.String())
	fmt.Printf("Block to write % X\n", blockData2)
	fmt.Printf("Sector %d\n", (blockAddr+1)/4)
	fmt.Printf("Block %d\n", blockAddr)
	ret, _, _ = syscall.SyscallN(uintptr(rfM1Write), 0, uintptr(blockAddr), uintptr(unsafe.Pointer(&blockData2)))
	fmt.Printf("Block to write as string: %s\n", data_to_write[:])
	handleFailedCallWithMessage(ret, rfM1WriteFuncName)
	fmt.Println("Block written", data_to_write[:])

}

func BruteForceAttack(rfM1Auth uintptr) []byte {
	var blockAddr uint32 = 0
	var ret uintptr

	// blockData2 data: 15F63DAF0A38 08 77 8F 69 17292304FFA7
	startKey := []byte{0x15, 0xE6, 0x3D, 0xAF, 0x0A, 0x38, 0x00, 0x00}
	endKey := []byte{0x15, 0xFF, 0x3D, 0xAF, 0x0A, 0x38, 0x00, 0x00}

	fmt.Printf("searching in range: % X - % X\n", startKey, endKey)
	fmt.Printf("start % X\n", startKey[7])
	fmt.Printf("end % X\n", endKey[7])

	start := binary.LittleEndian.Uint64(startKey)
	end := binary.LittleEndian.Uint64(endKey)

	fmt.Printf("start: %d\n", start)
	fmt.Printf("end: %d\n", end)

	for i := start; i <= end; i++ {
		key := make([]byte, 6)
		key_8 := make([]byte, 8)
		binary.LittleEndian.PutUint64(key_8, i)
		copy(key, key_8[0:6])

		fmt.Printf("Trying key: % X\n", key)

		ret, _, _ = syscall.SyscallN(uintptr(rfM1Auth), 0, 0x61, uintptr(blockAddr), uintptr(unsafe.Pointer(&key[0])))

		if ret == 0 {
			fmt.Println("Authentication successful!")
			fmt.Printf("Key found: % X\n", key)
			return key
		} else if ret != 22 {
			handleFailedCallWithMessage(ret, rfM1AuthFuncName)
		}
	}
	fmt.Printf("Key not found in range: % X - % X\n", startKey, endKey)
	os.Exit(1)
	return nil
}
