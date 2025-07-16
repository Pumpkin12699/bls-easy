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
	fmt.Println("=== BLS Easy 錯誤處理測試 ===")

	// 測試 1: 錯誤訊息功能
	fmt.Println("\n--- 測試 1: 錯誤訊息功能 ---")
	
	errorCodes := []C.int{
		C.BLS_SUCCESS,
		C.BLS_ERROR_INVALID_INPUT,
		C.BLS_ERROR_SIGNATURE_VERIFICATION_FAILED,
		C.BLS_ERROR_MEMORY_ALLOCATION_FAILED,
		C.BLS_ERROR_MEMORY_DEALLOCATION_FAILED,
		C.BLS_ERROR_SERIALIZATION_FAILED,
		C.BLS_ERROR_DESERIALIZATION_FAILED,
		C.BLS_ERROR_INVALID_KEY_FORMAT,
		C.BLS_ERROR_INVALID_SIGNATURE_FORMAT,
		C.BLS_ERROR_INVALID_MESSAGE_FORMAT,
		C.BLS_ERROR_KEY_GENERATION_FAILED,
		C.BLS_ERROR_INTERNAL_ERROR,
		-999, // 未知錯誤碼
	}

	for _, code := range errorCodes {
		message := C.bls_get_error_message(code)
		if message != nil {
			msgStr := C.GoString((*C.char)(message))
			fmt.Printf("錯誤碼 %d: %s\n", code, msgStr)
			C.bls_free_string(message)
		} else {
			fmt.Printf("錯誤碼 %d: 無法獲取錯誤訊息\n", code)
		}
	}

	// 測試 2: 錯誤碼檢查功能
	fmt.Println("\n--- 測試 2: 錯誤碼檢查功能 ---")
	
	testCodes := []C.int{0, -1, -5, -99, 1, 100}
	
	for _, code := range testCodes {
		isSuccess := C.bls_is_success(code)
		isError := C.bls_is_error(code)
		
		fmt.Printf("錯誤碼 %d: 成功=%d, 錯誤=%d\n", code, isSuccess, isError)
	}

	// 測試 3: 空指針錯誤處理
	fmt.Println("\n--- 測試 3: 空指針錯誤處理 ---")
	
	// 測試空私鑰
	result1 := C.bls_free_secret_key(nil)
	if result1 == C.BLS_ERROR_INVALID_INPUT {
		fmt.Println("✓ 空私鑰指針錯誤處理正確")
	} else {
		fmt.Printf("✗ 空私鑰指針錯誤處理異常，結果: %d\n", result1)
	}

	// 測試空公鑰
	result2 := C.bls_free_public_key(nil)
	if result2 == C.BLS_ERROR_INVALID_INPUT {
		fmt.Println("✓ 空公鑰指針錯誤處理正確")
	} else {
		fmt.Printf("✗ 空公鑰指針錯誤處理異常，結果: %d\n", result2)
	}

	// 測試空簽名
	result3 := C.bls_free_signature(nil)
	if result3 == C.BLS_ERROR_INVALID_INPUT {
		fmt.Println("✓ 空簽名指針錯誤處理正確")
	} else {
		fmt.Printf("✗ 空簽名指針錯誤處理異常，結果: %d\n", result3)
	}

	// 測試空字串
	result4 := C.bls_free_string(nil)
	if result4 == C.BLS_ERROR_INVALID_INPUT {
		fmt.Println("✓ 空字串指針錯誤處理正確")
	} else {
		fmt.Printf("✗ 空字串指針錯誤處理異常，結果: %d\n", result4)
	}

	// 測試 4: 簽名驗證錯誤處理
	fmt.Println("\n--- 測試 4: 簽名驗證錯誤處理 ---")
	
	// 產生有效的私鑰和公鑰
	secretKey := C.bls_generate_secret_key()
	if secretKey == nil {
		fmt.Println("✗ 無法產生私鑰")
		return
	}
	defer C.bls_free_secret_key(secretKey)

	publicKey := C.bls_get_public_key(secretKey)
	if publicKey == nil {
		fmt.Println("✗ 無法獲取公鑰")
		return
	}
	defer C.bls_free_public_key(publicKey)

	// 測試空簽名驗證
	result5 := C.bls_verify_signature(nil, C.CString("test"), publicKey)
	if result5 == C.BLS_ERROR_INVALID_INPUT {
		fmt.Println("✓ 空簽名驗證錯誤處理正確")
	} else {
		fmt.Printf("✗ 空簽名驗證錯誤處理異常，結果: %d\n", result5)
	}

	// 測試空訊息驗證
	result6 := C.bls_verify_signature(nil, nil, publicKey)
	if result6 == C.BLS_ERROR_INVALID_INPUT {
		fmt.Println("✓ 空訊息驗證錯誤處理正確")
	} else {
		fmt.Printf("✗ 空訊息驗證錯誤處理異常，結果: %d\n", result6)
	}

	// 測試空公鑰驗證
	result7 := C.bls_verify_signature(nil, C.CString("test"), nil)
	if result7 == C.BLS_ERROR_INVALID_INPUT {
		fmt.Println("✓ 空公鑰驗證錯誤處理正確")
	} else {
		fmt.Printf("✗ 空公鑰驗證錯誤處理異常，結果: %d\n", result7)
	}

	// 測試 5: 序列化錯誤處理
	fmt.Println("\n--- 測試 5: 序列化錯誤處理 ---")
	
	// 測試空公鑰序列化
	emptyPKStr := C.bls_public_key_to_string(nil)
	if emptyPKStr == nil {
		fmt.Println("✓ 空公鑰序列化錯誤處理正確")
	} else {
		fmt.Println("✗ 空公鑰序列化錯誤處理異常")
		C.bls_free_string(emptyPKStr)
	}

	// 測試空私鑰序列化
	emptySKStr := C.bls_secret_key_to_string(nil)
	if emptySKStr == nil {
		fmt.Println("✓ 空私鑰序列化錯誤處理正確")
	} else {
		fmt.Println("✗ 空私鑰序列化錯誤處理異常")
		C.bls_free_string(emptySKStr)
	}

	// 測試空簽名序列化
	emptySigStr := C.bls_signature_to_string(nil)
	if emptySigStr == nil {
		fmt.Println("✓ 空簽名序列化錯誤處理正確")
	} else {
		fmt.Println("✗ 空簽名序列化錯誤處理異常")
		C.bls_free_string(emptySigStr)
	}

	// 測試 6: 字串處理錯誤
	fmt.Println("\n--- 測試 6: 字串處理錯誤 ---")
	
	// 測試包含 null 字元的字串
	invalidMessage := C.CString("invalid\x00message")
	defer C.free(unsafe.Pointer(invalidMessage))
	
	signature := C.bls_sign(secretKey, invalidMessage)
	if signature == nil {
		fmt.Println("✓ 無效字串簽名錯誤處理正確")
	} else {
		fmt.Println("✗ 無效字串簽名錯誤處理異常")
		C.bls_free_signature(signature)
	}

	// 測試 7: 記憶體錯誤處理
	fmt.Println("\n--- 測試 7: 記憶體錯誤處理 ---")
	
	// 測試正常釋放
	validSK := C.bls_generate_secret_key()
	if validSK != nil {
		result8 := C.bls_free_secret_key(validSK)
		if result8 == C.BLS_SUCCESS {
			fmt.Println("✓ 正常釋放成功")
		} else {
			fmt.Printf("✗ 正常釋放失敗，錯誤碼: %d\n", result8)
		}
	}

	// 測試 8: 邊界條件
	fmt.Println("\n--- 測試 8: 邊界條件 ---")
	
	// 測試極長訊息
	longMessage := "This is a very long message that might cause issues " +
		"with memory allocation or string processing. " +
		"It contains many characters and might trigger edge cases " +
		"in the BLS signature implementation."
	
	cLongMessage := C.CString(longMessage)
	defer C.free(unsafe.Pointer(cLongMessage))
	
	longSignature := C.bls_sign(secretKey, cLongMessage)
	if longSignature != nil {
		fmt.Println("✓ 長訊息簽名成功")
		C.bls_free_signature(longSignature)
	} else {
		fmt.Println("✗ 長訊息簽名失敗")
	}

	// 測試空訊息
	emptyMessage := C.CString("")
	defer C.free(unsafe.Pointer(emptyMessage))
	
	emptySignature := C.bls_sign(secretKey, emptyMessage)
	if emptySignature != nil {
		fmt.Println("✓ 空訊息簽名成功")
		C.bls_free_signature(emptySignature)
	} else {
		fmt.Println("✗ 空訊息簽名失敗")
	}

	// 測試 9: 錯誤碼邊界測試
	fmt.Println("\n--- 測試 9: 錯誤碼邊界測試 ---")
	
	// 測試極大錯誤碼
	largeErrorCode := C.int(-999999)
	largeErrorMsg := C.bls_get_error_message(largeErrorCode)
	if largeErrorMsg != nil {
		msgStr := C.GoString((*C.char)(largeErrorMsg))
		fmt.Printf("極大錯誤碼 %d: %s\n", largeErrorCode, msgStr)
		C.bls_free_string(largeErrorMsg)
	}

	// 測試極小錯誤碼
	smallErrorCode := C.int(999999)
	smallErrorMsg := C.bls_get_error_message(smallErrorCode)
	if smallErrorMsg != nil {
		msgStr := C.GoString((*C.char)(smallErrorMsg))
		fmt.Printf("極小錯誤碼 %d: %s\n", smallErrorCode, msgStr)
		C.bls_free_string(smallErrorMsg)
	}

	fmt.Println("\n=== 錯誤處理測試完成 ===")
} 