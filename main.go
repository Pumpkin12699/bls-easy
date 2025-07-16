package main

// #cgo LDFLAGS: -L./target/release -lw3f_bls_easy
// #include <stdlib.h>
// #include "./src/bls_easy_ffi.h"
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println("=== BLS Easy FFI 基本功能測試 ===")

	// 1. 產生私鑰
	fmt.Println("1. 產生私鑰...")
	secretKey := C.bls_generate_secret_key()
	if secretKey == nil {
		fmt.Println("Failed to generate secret key")
		return
	}
	defer C.bls_free_secret_key(secretKey)
	fmt.Println("✓ 私鑰產生成功")

	// 2. 匯出公鑰
	fmt.Println("2. 匯出公鑰...")
	publicKey := C.bls_get_public_key(secretKey)
	if publicKey == nil {
		fmt.Println("Failed to get public key")
		return
	}
	defer C.bls_free_public_key(publicKey)
	fmt.Println("✓ 公鑰匯出成功")

	// 3. 簽署訊息
	fmt.Println("3. 簽署訊息...")
	message := "Hello, BLS Easy World!"
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))

	signature := C.bls_sign(secretKey, cMessage)
	if signature == nil {
		fmt.Println("Failed to sign message")
		return
	}
	defer C.bls_free_signature(signature)
	fmt.Println("✓ 簽名產生成功")

	// 4. 驗證簽名
	fmt.Println("4. 驗證簽名...")
	result := C.bls_verify_signature(signature, cMessage, publicKey)
	if result == C.BLS_SUCCESS {
		fmt.Println("✓ 簽名驗證成功!")
	} else {
		fmt.Printf("✗ 簽名驗證失敗，錯誤碼: %d\n", result)
		return
	}

	// 5. 序列化測試
	fmt.Println("5. 序列化測試...")
	secretKeyStr := C.bls_secret_key_to_string(secretKey)
	defer C.bls_free_string(secretKeyStr)
	secretKeyGoStr := C.GoString((*C.char)(secretKeyStr))
	fmt.Printf("私鑰: %s\n", secretKeyGoStr)

	publicKeyStr := C.bls_public_key_to_string(publicKey)
	defer C.bls_free_string(publicKeyStr)
	publicKeyGoStr := C.GoString((*C.char)(publicKeyStr))
	fmt.Printf("公鑰: %s\n", publicKeyGoStr)

	signatureStr := C.bls_signature_to_string(signature)
	defer C.bls_free_string(signatureStr)
	signatureGoStr := C.GoString((*C.char)(signatureStr))
	fmt.Printf("簽名: %s\n", signatureGoStr)

	// 6. 錯誤訊息測試
	fmt.Println("6. 錯誤處理測試...")
	badMessage := "Different message"
	cBadMessage := C.CString(badMessage)
	defer C.free(unsafe.Pointer(cBadMessage))

	result2 := C.bls_verify_signature(signature, cBadMessage, publicKey)
	if result2 == C.BLS_ERROR_SIGNATURE_VERIFICATION_FAILED {
		fmt.Println("✓ 錯誤訊息正確被拒絕")
	} else {
		fmt.Printf("✗ 錯誤訊息處理異常，結果: %d\n", result2)
	}

	fmt.Println("\n=== 所有基本功能測試完成 ===")
} 