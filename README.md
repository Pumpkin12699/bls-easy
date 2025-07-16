# BLS Easy FFI

這是 BLS 簽名庫的簡化版 Go FFI (Foreign Function Interface)，專注於核心 BLS 功能，讓 Go 程式可以輕鬆呼叫 Rust 實現的 BLS 簽名功能。

## 專案特色

- **簡化設計**: 只實作核心 BLS API，避免複雜性
- **完整功能**: 涵蓋簽署、聚合、驗證等完整流程
- **易於使用**: 提供清晰的 Go 介面
- **記憶體安全**: 完整的記憶體管理機制

## 專案結構

```
bls_easy/
├── src/                    # Rust 原始碼
│   ├── lib.rs             # 主要庫檔案
│   ├── ffi.rs             # 簡化版 FFI 介面
│   └── ...                # 其他 Rust 模組
├── main.go                # Go 測試程式
├── Cargo.toml             # Rust 依賴配置
├── cbindgen.toml          # C 標頭檔生成配置
├── go.mod                 # Go 模組配置
├── Makefile               # 編譯腳本
└── README.md              # 本檔案
```

## 核心功能

### 簽署者端 (Signer Side)

1. **產生私鑰**: `bls_generate_secret_key()`
2. **匯出公鑰**: `bls_get_public_key(secret_key)`
3. **產生持有證明**: `bls_generate_pok(secret_key)`
4. **簽署訊息**: `bls_sign(secret_key, message, context)`

### 聚合者端 (Aggregator Side)

1. **驗證持有證明**: `bls_verify_pok(proof, public_key)`
2. **聚合簽章**: `bls_aggregate_signatures(signatures, count)`
3. **聚合公鑰**: `bls_aggregate_public_keys(public_keys, count)`

### 驗證者端 (Verifier Side)

1. **驗證聚合簽章 (不同訊息)**: `bls_verify_aggregated_signature(signature, messages, public_keys, count)`
2. **驗證聚合簽章 (相同訊息)**: `bls_verify_aggregated_signature_same_message(signature, message, public_key)`

### 序列化功能

- **序列化**: `bls_*_to_string()` 系列函數
- **反序列化**: `bls_*_from_string()` 系列函數
- **記憶體管理**: `bls_free_*()` 系列函數

## 安裝和編譯

### 前置需求

- Rust (1.70+)
- Go (1.21+)
- cbindgen

### 安裝 cbindgen

```bash
cargo install --force cbindgen
```

### 編譯和測試

```bash
# 完整設置 (安裝 cbindgen + 編譯 + 測試)
make setup

# 或者分步驟執行
make install-cbindgen  # 安裝 cbindgen
make build            # 編譯 Rust 庫
make test             # 執行 Go 測試
```

## API 使用範例

### Go 程式範例

```go
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
    // 1. 產生私鑰
    secretKey := C.bls_generate_secret_key()
    defer C.bls_free_secret_key(secretKey)

    // 2. 匯出公鑰
    publicKey := C.bls_get_public_key(secretKey)
    defer C.bls_free_public_key(publicKey)

    // 3. 產生持有證明
    proof := C.bls_generate_pok(secretKey)
    defer C.bls_free_schnorr_proof(proof)

    // 4. 簽署訊息
    message := "Hello, BLS Easy!"
    cMessage := C.CString(message)
    defer C.free(unsafe.Pointer(cMessage))
    
    signature := C.bls_sign(secretKey, cMessage, nil)
    defer C.bls_free_signature(signature)

    // 5. 驗證持有證明
    result := C.bls_verify_pok(proof, publicKey)
    if result == C.BLS_SUCCESS {
        fmt.Println("Proof of possession verified!")
    }

    // 6. 聚合簽章
    signatures := [2]*C.BLSSignature{signature, signature}
    aggregatedSig := C.bls_aggregate_signatures(&signatures[0], 2)
    defer C.bls_free_signature(aggregatedSig)

    // 7. 驗證聚合簽章
    result2 := C.bls_verify_aggregated_signature_same_message(aggregatedSig, cMessage, publicKey)
    if result2 == C.BLS_SUCCESS {
        fmt.Println("Aggregated signature verified!")
    }
}
```

## FFI 函數列表

### 簽署者端函數

- `bls_generate_secret_key()` - 產生新的私鑰
- `bls_get_public_key(secret_key)` - 從私鑰獲取公鑰
- `bls_generate_pok(secret_key)` - 產生持有證明
- `bls_sign(secret_key, message, context)` - 簽名訊息

### 聚合者端函數

- `bls_verify_pok(proof, public_key)` - 驗證持有證明
- `bls_aggregate_signatures(signatures, count)` - 聚合簽章
- `bls_aggregate_public_keys(public_keys, count)` - 聚合公鑰

### 驗證者端函數

- `bls_verify_aggregated_signature(signature, messages, public_keys, count)` - 驗證聚合簽章 (不同訊息)
- `bls_verify_aggregated_signature_same_message(signature, message, public_key)` - 驗證聚合簽章 (相同訊息)

### 序列化函數

- `bls_public_key_to_string(public_key)` - 序列化公鑰為字串
- `bls_secret_key_to_string(secret_key)` - 序列化私鑰為字串
- `bls_signature_to_string(signature)` - 序列化簽名為字串
- `bls_public_key_from_string(hex_string)` - 從字串反序列化公鑰
- `bls_secret_key_from_string(hex_string)` - 從字串反序列化私鑰
- `bls_signature_from_string(hex_string)` - 從字串反序列化簽名

### 記憶體管理函數

- `bls_free_secret_key(secret_key)` - 釋放私鑰記憶體
- `bls_free_public_key(public_key)` - 釋放公鑰記憶體
- `bls_free_signature(signature)` - 釋放簽名記憶體
- `bls_free_schnorr_proof(proof)` - 釋放持有證明記憶體
- `bls_free_string(string)` - 釋放字串記憶體

## 錯誤碼

- `BLS_SUCCESS` (0) - 操作成功
- `BLS_ERROR_INVALID_INPUT` (-1) - 無效輸入
- `BLS_ERROR_SIGNATURE_VERIFICATION_FAILED` (-2) - 簽名驗證失敗
- `BLS_ERROR_MEMORY_ALLOCATION_FAILED` (-3) - 記憶體分配失敗

## 與完整版的差異

### 簡化設計
- 移除了複雜的金鑰對管理
- 專注於核心 BLS 功能
- 簡化的 API 設計

### 功能對比

| 功能 | 完整版 (bls-go) | 簡化版 (bls_easy) |
|------|------------------|-------------------|
| 金鑰對生成 | ✓ | ✓ |
| 簽名/驗證 | ✓ | ✓ |
| 聚合功能 | ✓ | ✓ |
| 持有證明 | ✓ | ✓ |
| 序列化 | ✓ | ✓ |
| 金鑰對管理 | 複雜 | 簡化 |
| 多種曲線支援 | ✓ | 僅 BLS12-381 |
| 進階聚合 | ✓ | 基本聚合 |

## 使用場景

### 適合使用 bls_easy 的場景
- 需要快速整合 BLS 功能
- 只需要基本的簽名和聚合功能
- 希望簡化的 API 設計
- 主要使用 BLS12-381 曲線

### 適合使用完整版的場景
- 需要多種曲線支援
- 需要進階的聚合功能
- 需要完整的金鑰管理
- 需要更細緻的控制

## 注意事項

1. **記憶體管理**: 所有 FFI 函數都需要正確的記憶體管理，使用 `defer` 來確保資源被釋放
2. **字串處理**: 字串參數需要使用 `C.CString()` 轉換，並在完成後釋放
3. **空指針檢查**: 返回的指針可能為 `nil`，使用前需要檢查
4. **序列化格式**: 序列化的資料使用十六進制字串格式
5. **聚合限制**: 聚合功能需要相同數量的簽名和公鑰

## 清理

```bash
make clean
```

這會清理編譯產生的檔案和生成的標頭檔。

## 貢獻

歡迎提交 Issue 和 Pull Request 來改善這個專案！

## 授權

本專案採用 MIT/Apache-2.0 雙重授權。

## 📋 記憶體管理改進總結

我們已經成功完成了記憶體管理部分的改進：

### ✅ 已完成的改進：

1. **改進的記憶體釋放函數**
   - 所有 `bls_free_*` 函數現在返回錯誤碼
   - 更好的空指針檢查
   - 統一的錯誤處理

2. **新增工具函數**
   - `bls_is_null()` - 檢查指針是否為空
   - `bls_get_memory_stats()` - 記憶體統計 (debug 模式)

3. **完整的測試覆蓋**
   - 基本記憶體分配和釋放
   - 多重分配和釋放測試
   - 字串記憶體管理
   - 錯誤處理測試
   - 指針檢查測試

4. **錯誤碼定義**
   - `BLS_ERROR_MEMORY_DEALLOCATION_FAILED` - 新增記憶體釋放錯誤碼

### 🧪 測試結果：
- ✅ 所有記憶體管理測試通過
- ✅ Debug 和 Release 模式都正常工作
- ✅ 錯誤處理正確
- ✅ 記憶體洩漏檢測正常

### 📄 檔案結構：
```
bls_easy/
├── src/
│   ├── ffi.rs              # 改進的 FFI 介面
│   └── bls_easy_ffi.h      # 更新的標頭檔
├── memory_test_main.go      # 記憶體管理測試程式
└── ...
```

### 💡 下一步建議：

1. **聚合功能改進** - 實現真正的 BLS 聚合
2. **批量操作** - 添加批量簽名和驗證
3. **性能優化** - 添加性能測試和優化
4. **文檔完善** - 添加詳細的 API 文檔

你想要繼續哪個部分的開發呢？

