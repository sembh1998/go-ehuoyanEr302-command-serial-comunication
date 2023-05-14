package main

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

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
		syscall.SyscallN(rfClosePort)
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
	time.Sleep(1 * time.Second)

	// Call the rf_light function
	ret, _, _ = syscall.SyscallN(uintptr(rfLight), 0, 2)
	handleFailedCallWithMessage(ret, rfLightFuncName)
	fmt.Println("LED turned red")
	// Wait for 1 second
	time.Sleep(1 * time.Second)

	// Call the rf_light function
	ret, _, _ = syscall.SyscallN(uintptr(rfLight), 0, 1)
	handleFailedCallWithMessage(ret, rfLightFuncName)
	fmt.Println("LED turned blue")
	// Wait for 1 second
	time.Sleep(1 * time.Second)

	// Call the rf_beep function
	ret, _, _ = syscall.SyscallN(uintptr(rfBeep), 0, 10)
	handleFailedCallWithMessage(ret, rfBeepFuncName)
	fmt.Println("Beeped")
	time.Sleep(1 * time.Second)

	// Call the rf_beep function
	ret, _, _ = syscall.SyscallN(uintptr(rfBeep), 0, 10)
	handleFailedCallWithMessage(ret, rfBeepFuncName)
	fmt.Println("Beeped")
	time.Sleep(1 * time.Second)

	// Call the rf_request function
	var tagType uint32
	ret, _, _ = syscall.SyscallN(uintptr(rfRequest), 0, 1, uintptr(unsafe.Pointer(&tagType)))
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
	key := [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	ret, _, _ = syscall.SyscallN(uintptr(rfM1Auth), 0, 0x60, 0, uintptr(unsafe.Pointer(&key)))
	handleFailedCallWithMessage(ret, rfM1AuthFuncName)
	fmt.Println("Authenticated")

	// Call the rf_M1_read function
	var blockData [16]byte
	var blockLen uint32
	var blockAddr uint32 = 1
	ret, _, _ = syscall.SyscallN(uintptr(rfM1Read), 0, uintptr(blockAddr), uintptr(unsafe.Pointer(&blockData)), uintptr(unsafe.Pointer(&blockLen)))
	handleFailedCallWithMessage(ret, rfM1ReadFuncName)
	fmt.Printf("Block data: % X\n", blockData)
	fmt.Printf("Block length: %d\n", blockLen)
	// print the block data as a string
	fmt.Printf("Block data as string: %s\n", blockData[:])

	// Call the rf_M1_write function
	// save random data to block 1
	var blockData2 [16]byte
	var data = []byte(faker.Commerce().ProductName())
	copy(blockData2[:], data)
	fmt.Println("Block to write", blockData2[:])
	fmt.Printf("Block to write as string: %s\n", blockData2[:])
	ret, _, _ = syscall.SyscallN(uintptr(rfM1Write), 0, uintptr(blockAddr), uintptr(unsafe.Pointer(&blockData2)))
	handleFailedCallWithMessage(ret, rfM1WriteFuncName)
	fmt.Println("Block written", blockData2[:])

}
