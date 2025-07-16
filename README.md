# bls-easy

BLS Easy 是一個以 Rust 實作、支援 Go 調用的 BLS 簽章 FFI 函式庫，專注於簡單、健壯、易於整合。

---

## 目前進度與功能

### 1. Rust BLS FFI Library 基礎建設
- 以 Rust 撰寫 BLS 簽章核心邏輯，並以 C ABI 封裝，支援 Go 調用。
- 提供私鑰產生、公鑰導出、簽名、驗證、序列化等基本功能。

### 2. 記憶體管理
- 所有 FFI 產生的物件（私鑰、公鑰、簽名、字串）皆有對應的釋放函數。
- 釋放函數皆有錯誤碼回傳，並對空指針做防呆。
- 提供指針檢查、記憶體統計（debug 模式）等輔助函數。
- Go 端有完整的記憶體管理測試程式（`memory_test_main.go`）。

### 3. 錯誤處理系統
- 定義多種錯誤碼（如輸入錯誤、序列化錯誤、記憶體錯誤、簽名驗證失敗等）。
- 提供 `bls_get_error_message` 查詢錯誤碼對應訊息。
- 提供 `bls_is_success`、`bls_is_error` 判斷 API。
- 所有 API 都有嚴謹的錯誤處理與回傳。
- Go 端有完整的錯誤處理測試程式（`error_test_main.go`），涵蓋各種情境與邊界條件。

### 4. 測試與驗證
- `memory_test_main.go`：測試記憶體分配、釋放、字串管理、錯誤處理。
- `error_test_main.go`：測試所有錯誤碼、錯誤訊息、API 錯誤行為與邊界條件。
- 測試皆通過，確保 FFI 介面健壯且安全。

### 5. 文件與備份
- `BACKUP.md`：保留重要狀態的備份說明。

---

## 主要 API 說明

### 產生私鑰
```c
BLSSecretKey *bls_generate_secret_key(void);
```

### 取得公鑰
```c
BLSPublicKey *bls_get_public_key(const BLSSecretKey *secret_key);
```

### 簽名
```c
BLSSignature *bls_sign(BLSSecretKey *secret_key, const char *message);
```

### 驗證簽名
```c
int bls_verify_signature(const BLSSignature *signature, const char *message, const BLSPublicKey *public_key);
```

### 釋放記憶體
```c
int bls_free_secret_key(BLSSecretKey *secret_key);
int bls_free_public_key(BLSPublicKey *public_key);
int bls_free_signature(BLSSignature *signature);
int bls_free_string(char *s);
```

### 錯誤處理
```c
char *bls_get_error_message(int error_code);
int bls_is_success(int error_code);
int bls_is_error(int error_code);
```

### 其他輔助
```c
int bls_is_null(const void *ptr);
char *bls_get_memory_stats(void); // 僅 debug 模式
```

---

## 測試方式

### 記憶體管理測試
```sh
go run memory_test_main.go
```

### 錯誤處理測試
```sh
go run error_test_main.go
```

---

## 下一步建議
- 支援 BLS 聚合簽名與批次驗證
- 增加更多序列化/反序列化 API
- 完善 API 文件與範例

---

如有任何問題或建議，歡迎開 issue 或 PR！

