package main

// #cgo LDFLAGS: -L./target/debug -lw3f_bls_easy
// #include <stdlib.h>
// #include "./src/bls_easy_ffi.h"
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println("=== BLS Easy 記憶體管理測試 (Debug 模式) ===")

	// 測試 1: 基本記憶體分配和釋放
	fmt.Println("\n--- 測試 1: 基本記憶體分配和釋放 ---")
	
	// 產生私鑰
	secretKey := C.bls_generate_secret_key()
	if secretKey == nil {
		fmt.Println("Failed to generate secret key")
		return
	}
	fmt.Println("✓ 私鑰分配成功")

	// 檢查指針是否為空
	isNull := C.bls_is_null(unsafe.Pointer(secretKey))
	if isNull == 1 {
		fmt.Println("✗ 指針檢查失敗 - 指針被誤判為空")
		return
	} else {
		fmt.Println("✓ 指針檢查通過")
	}

	// 釋放私鑰
	result := C.bls_free_secret_key(secretKey)
	if result == C.BLS_SUCCESS {
		fmt.Println("✓ 私鑰釋放成功")
	} else {
		fmt.Printf("✗ 私鑰釋放失敗，錯誤碼: %d\n", result)
	}

	// 測試 2: 多重分配和釋放
	fmt.Println("\n--- 測試 2: 多重分配和釋放 ---")
	
	var secretKeys []*C.BLSSecretKey
	var publicKeys []*C.BLSPublicKey
	var signatures []*C.BLSSignature

	// 分配多個物件
	for i := 0; i < 5; i++ {
		sk := C.bls_generate_secret_key()
		if sk == nil {
			fmt.Printf("✗ 第 %d 個私鑰分配失敗\n", i+1)
			continue
		}
		secretKeys = append(secretKeys, sk)

		pk := C.bls_get_public_key(sk)
		if pk == nil {
			fmt.Printf("✗ 第 %d 個公鑰分配失敗\n", i+1)
			continue
		}
		publicKeys = append(publicKeys, pk)

		message := fmt.Sprintf("Test message %d", i+1)
		cMessage := C.CString(message)
		defer C.free(unsafe.Pointer(cMessage))

		sig := C.bls_sign(sk, cMessage)
		if sig == nil {
			fmt.Printf("✗ 第 %d 個簽名分配失敗\n", i+1)
			continue
		}
		signatures = append(signatures, sig)

		fmt.Printf("✓ 第 %d 組物件分配成功\n", i+1)
	}

	// 釋放所有物件
	fmt.Println("開始釋放物件...")
	
	for i, sig := range signatures {
		result := C.bls_free_signature(sig)
		if result != C.BLS_SUCCESS {
			fmt.Printf("✗ 第 %d 個簽名釋放失敗，錯誤碼: %d\n", i+1, result)
		}
	}

	for i, pk := range publicKeys {
		result := C.bls_free_public_key(pk)
		if result != C.BLS_SUCCESS {
			fmt.Printf("✗ 第 %d 個公鑰釋放失敗，錯誤碼: %d\n", i+1, result)
		}
	}

	for i, sk := range secretKeys {
		result := C.bls_free_secret_key(sk)
		if result != C.BLS_SUCCESS {
			fmt.Printf("✗ 第 %d 個私鑰釋放失敗，錯誤碼: %d\n", i+1, result)
		}
	}

	fmt.Println("✓ 所有物件釋放完成")

	// 測試 3: 字串記憶體管理
	fmt.Println("\n--- 測試 3: 字串記憶體管理 ---")
	
	// 產生測試物件
	testSecretKey := C.bls_generate_secret_key()
	defer C.bls_free_secret_key(testSecretKey)

	testPublicKey := C.bls_get_public_key(testSecretKey)
	defer C.bls_free_public_key(testPublicKey)

	// 序列化為字串
	secretKeyStr := C.bls_secret_key_to_string(testSecretKey)
	if secretKeyStr == nil {
		fmt.Println("✗ 私鑰序列化失敗")
		return
	}
	fmt.Println("✓ 私鑰序列化成功")

	publicKeyStr := C.bls_public_key_to_string(testPublicKey)
	if publicKeyStr == nil {
		fmt.Println("✗ 公鑰序列化失敗")
		return
	}
	fmt.Println("✓ 公鑰序列化成功")

	// 檢查字串內容
	secretKeyGoStr := C.GoString((*C.char)(secretKeyStr))
	publicKeyGoStr := C.GoString((*C.char)(publicKeyStr))

	fmt.Printf("私鑰字串長度: %d\n", len(secretKeyGoStr))
	fmt.Printf("公鑰字串長度: %d\n", len(publicKeyGoStr))

	// 釋放字串記憶體
	result1 := C.bls_free_string(secretKeyStr)
	if result1 == C.BLS_SUCCESS {
		fmt.Println("✓ 私鑰字串釋放成功")
	} else {
		fmt.Printf("✗ 私鑰字串釋放失敗，錯誤碼: %d\n", result1)
	}

	result2 := C.bls_free_string(publicKeyStr)
	if result2 == C.BLS_SUCCESS {
		fmt.Println("✓ 公鑰字串釋放成功")
	} else {
		fmt.Printf("✗ 公鑰字串釋放失敗，錯誤碼: %d\n", result2)
	}

	// 測試 4: 錯誤處理
	fmt.Println("\n--- 測試 4: 錯誤處理 ---")
	
	// 測試空指針釋放
	result3 := C.bls_free_secret_key(nil)
	if result3 == C.BLS_ERROR_INVALID_INPUT {
		fmt.Println("✓ 空指針錯誤處理正確")
	} else {
		fmt.Printf("✗ 空指針錯誤處理異常，結果: %d\n", result3)
	}

	// 測試空字串釋放
	result4 := C.bls_free_string(nil)
	if result4 == C.BLS_ERROR_INVALID_INPUT {
		fmt.Println("✓ 空字串錯誤處理正確")
	} else {
		fmt.Printf("✗ 空字串錯誤處理異常，結果: %d\n", result4)
	}

	// 測試指針檢查
	isNull1 := C.bls_is_null(nil)
	if isNull1 == 1 {
		fmt.Println("✓ 空指針檢查正確")
	} else {
		fmt.Println("✗ 空指針檢查異常")
	}

	isNull2 := C.bls_is_null(unsafe.Pointer(testSecretKey))
	if isNull2 == 0 {
		fmt.Println("✓ 非空指針檢查正確")
	} else {
		fmt.Println("✗ 非空指針檢查異常")
	}

	// 測試 5: 記憶體統計 (Debug 模式)
	fmt.Println("\n--- 測試 5: 記憶體統計 (Debug 模式) ---")
	
	stats := C.bls_get_memory_stats()
	if stats != nil {
		statsStr := C.GoString((*C.char)(stats))
		fmt.Printf("記憶體統計: %s\n", statsStr)
		C.bls_free_string(stats)
	} else {
		fmt.Println("記憶體統計功能不可用")
	}

	fmt.Println("\n=== 記憶體管理測試完成 ===")
} 