package gateways

import (
	"fmt"
	"log"
	"syscall"
)

type ledColor int
type LEDColor ledColor

const (
	libPath             = "MasterRD.dll"
	libVerFuncName      = "lib_ver"
	rfInitComFuncName   = "rf_init_com"
	rfClosePortFuncName = "rf_ClosePort"
	rfLightFuncName     = "rf_light"
	rfBeepFuncName      = "rf_beep"
	rfRequestFuncName   = "rf_request"
	rfAnticollFuncName  = "rf_anticoll"
	rfSelectFuncName    = "rf_select"
	rfM1AuthFuncName    = "rf_M1_authentication2"
	rfM1ReadFuncName    = "rf_M1_read"
	rfM1WriteFuncName   = "rf_M1_write"

	// Colors
	LedOff  ledColor = 0
	LedBlue ledColor = 1
	LedRed  ledColor = 2
)

var (
	handle      syscall.Handle
	rfLibVer    uintptr
	rfInitCom   uintptr
	rfClosePort uintptr
	rfLight     uintptr
	rfBeep      uintptr
	rfRequest   uintptr
	rfAnticoll  uintptr
	rfSelect    uintptr
	rfM1Auth    uintptr
	rfM1Read    uintptr
	rfM1Write   uintptr
)

func init() {
	handle, err := syscall.LoadLibrary(libPath)
	if err != nil {
		log.Fatalf("Failed to load DLL: %v", err)
	}

	getProcAddress := func(name string) uintptr {
		proc, err := syscall.GetProcAddress(handle, name)
		if err != nil {
			log.Fatalf("Failed to get function address: %s: %v", name, err)
		}
		return proc
	}

	rfLibVer = getProcAddress(libVerFuncName)
	rfInitCom = getProcAddress(rfInitComFuncName)
	rfClosePort = getProcAddress(rfClosePortFuncName)
	rfLight = getProcAddress(rfLightFuncName)
	rfBeep = getProcAddress(rfBeepFuncName)
	rfRequest = getProcAddress(rfRequestFuncName)
	rfAnticoll = getProcAddress(rfAnticollFuncName)
	rfSelect = getProcAddress(rfSelectFuncName)
	rfM1Auth = getProcAddress(rfM1AuthFuncName)
	rfM1Read = getProcAddress(rfM1ReadFuncName)
	rfM1Write = getProcAddress(rfM1WriteFuncName)
}

func handleError(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func handleFailedCallWithMessage(ret uintptr, msg string) {
	if ret != 0 {
		ret2, _, _ := syscall.SyscallN(rfClosePort)
		log.Printf("rf_ClosePort returned %d", ret2)
		log.Fatalf("%s failed with error code %d", msg, ret)
	}
}

func InitializeCardReader() {
	// Call the rf_init_com function
	ret, _, _ := syscall.SyscallN(uintptr(rfInitCom), 3, 115200)
	handleFailedCallWithMessage(ret, rfInitComFuncName)
	fmt.Println("Connected successfully")
}

func LightLED(color ledColor) {
	ret, _, _ := syscall.SyscallN(uintptr(rfLight), 0, uintptr(color))
	handleFailedCallWithMessage(ret, rfLightFuncName)
	if color == LedOff {
		fmt.Println("LED turned off")
	} else if color == LedBlue {
		fmt.Println("LED turned blue")
	} else if color == LedRed {
		fmt.Println("LED turned red")
	}
}
