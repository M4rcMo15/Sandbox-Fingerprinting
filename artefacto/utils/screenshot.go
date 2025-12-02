package utils

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"unsafe"

	"golang.org/x/sys/windows"
)

// CaptureScreenshot captura la pantalla y la devuelve como base64
func CaptureScreenshot() string {
	user32 := windows.NewLazyDLL("user32.dll")
	gdi32 := windows.NewLazyDLL("gdi32.dll")

	getDC := user32.NewProc("GetDC")
	createCompatibleDC := gdi32.NewProc("CreateCompatibleDC")
	createCompatibleBitmap := gdi32.NewProc("CreateCompatibleBitmap")
	selectObject := gdi32.NewProc("SelectObject")
	bitBlt := gdi32.NewProc("BitBlt")
	getSystemMetrics := user32.NewProc("GetSystemMetrics")
	getDIBits := gdi32.NewProc("GetDIBits")
	deleteDC := gdi32.NewProc("DeleteDC")
	deleteObject := gdi32.NewProc("DeleteObject")
	releaseDC := user32.NewProc("ReleaseDC")

	// Obtener dimensiones de la pantalla
	screenWidth, _, _ := getSystemMetrics.Call(0)  // SM_CXSCREEN
	screenHeight, _, _ := getSystemMetrics.Call(1) // SM_CYSCREEN

	if screenWidth == 0 || screenHeight == 0 {
		return ""
	}

	// Obtener DC de la pantalla
	hdcScreen, _, _ := getDC.Call(0)
	if hdcScreen == 0 {
		return ""
	}
	defer releaseDC.Call(0, hdcScreen)

	// Crear DC compatible
	hdcMem, _, _ := createCompatibleDC.Call(hdcScreen)
	if hdcMem == 0 {
		return ""
	}
	defer deleteDC.Call(hdcMem)

	// Crear bitmap compatible
	hBitmap, _, _ := createCompatibleBitmap.Call(hdcScreen, screenWidth, screenHeight)
	if hBitmap == 0 {
		return ""
	}
	defer deleteObject.Call(hBitmap)

	// Seleccionar bitmap en el DC
	selectObject.Call(hdcMem, hBitmap)

	// Copiar la pantalla al bitmap
	const SRCCOPY = 0x00CC0020
	bitBlt.Call(hdcMem, 0, 0, screenWidth, screenHeight, hdcScreen, 0, 0, SRCCOPY)

	// Preparar estructura BITMAPINFO
	type BITMAPINFOHEADER struct {
		BiSize          uint32
		BiWidth         int32
		BiHeight        int32
		BiPlanes        uint16
		BiBitCount      uint16
		BiCompression   uint32
		BiSizeImage     uint32
		BiXPelsPerMeter int32
		BiYPelsPerMeter int32
		BiClrUsed       uint32
		BiClrImportant  uint32
	}

	type BITMAPINFO struct {
		BmiHeader BITMAPINFOHEADER
		BmiColors [1]uint32
	}

	bi := BITMAPINFO{}
	bi.BmiHeader.BiSize = uint32(unsafe.Sizeof(bi.BmiHeader))
	bi.BmiHeader.BiWidth = int32(screenWidth)
	bi.BmiHeader.BiHeight = -int32(screenHeight) // Top-down
	bi.BmiHeader.BiPlanes = 1
	bi.BmiHeader.BiBitCount = 32
	bi.BmiHeader.BiCompression = 0 // BI_RGB

	// Calcular tamaño del buffer
	bufferSize := int(screenWidth) * int(screenHeight) * 4
	buffer := make([]byte, bufferSize)

	// Obtener bits del bitmap
	ret, _, _ := getDIBits.Call(
		hdcMem,
		hBitmap,
		0,
		uintptr(screenHeight),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&bi)),
		0, // DIB_RGB_COLORS
	)

	if ret == 0 {
		return ""
	}

	// Convertir a imagen Go
	img := image.NewRGBA(image.Rect(0, 0, int(screenWidth), int(screenHeight)))
	
	for y := 0; y < int(screenHeight); y++ {
		for x := 0; x < int(screenWidth); x++ {
			offset := (y*int(screenWidth) + x) * 4
			// BGRA -> RGBA
			img.Pix[offset+0] = buffer[offset+2]
			img.Pix[offset+1] = buffer[offset+1]
			img.Pix[offset+2] = buffer[offset+0]
			img.Pix[offset+3] = buffer[offset+3]
		}
	}

	// Codificar a PNG
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return ""
	}
	encoded_bytes := base64.StdEncoding.EncodeToString(buf.Bytes())
	return encoded_bytes
}