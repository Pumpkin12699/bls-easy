use std::os::raw::{c_char, c_int};
use std::ffi::{CStr, CString};
use std::ptr;

use crate::single::{Keypair, PublicKey, SecretKey, Signature};
use crate::engine::ZBLS;
use crate::Message;
use crate::serialize::SerializableToBytes;

// 定義錯誤碼
pub const BLS_SUCCESS: c_int = 0;
pub const BLS_ERROR_INVALID_INPUT: c_int = -1;
pub const BLS_ERROR_SIGNATURE_VERIFICATION_FAILED: c_int = -2;
pub const BLS_ERROR_MEMORY_ALLOCATION_FAILED: c_int = -3;
pub const BLS_ERROR_MEMORY_DEALLOCATION_FAILED: c_int = -4;

// 定義結構體來包裝 Rust 類型
#[repr(C)]
pub struct BLSKeypair {
    keypair: Keypair<ZBLS>,
}

#[repr(C)]
pub struct BLSPublicKey {
    public_key: PublicKey<ZBLS>,
}

#[repr(C)]
pub struct BLSSecretKey {
    secret_key: SecretKey<ZBLS>,
}

#[repr(C)]
pub struct BLSSignature {
    signature: Signature<ZBLS>,
}

// ===== 基本功能測試 =====

// 1. 產生私鑰
#[no_mangle]
pub extern "C" fn bls_generate_secret_key() -> *mut BLSSecretKey {
    use rand::thread_rng;
    let secret_key = SecretKey::<ZBLS>::generate(thread_rng());
    let boxed = Box::new(BLSSecretKey { secret_key });
    Box::into_raw(boxed)
}

// 2. 匯出公鑰
#[no_mangle]
pub extern "C" fn bls_get_public_key(secret_key: *const BLSSecretKey) -> *mut BLSPublicKey {
    if secret_key.is_null() {
        return ptr::null_mut();
    }
    
    let sk = unsafe { &*secret_key };
    let public_key = sk.secret_key.into_public();
    let boxed = Box::new(BLSPublicKey { public_key });
    Box::into_raw(boxed)
}

// 3. 簽署訊息
#[no_mangle]
pub extern "C" fn bls_sign(
    secret_key: *mut BLSSecretKey,
    message: *const c_char,
) -> *mut BLSSignature {
    if secret_key.is_null() || message.is_null() {
        return ptr::null_mut();
    }
    
    let sk = unsafe { &mut *secret_key };
    let msg_str = unsafe { CStr::from_ptr(message) };
    let msg_bytes = msg_str.to_bytes();
    
    let message_obj = Message::new(b"", msg_bytes);
    
    use rand::thread_rng;
    let signature = sk.secret_key.sign(&message_obj, thread_rng());
    let boxed = Box::new(BLSSignature { signature });
    Box::into_raw(boxed)
}

// 4. 驗證簽名
#[no_mangle]
pub extern "C" fn bls_verify_signature(
    signature: *const BLSSignature,
    message: *const c_char,
    public_key: *const BLSPublicKey,
) -> c_int {
    if signature.is_null() || message.is_null() || public_key.is_null() {
        return BLS_ERROR_INVALID_INPUT;
    }
    
    let sig = unsafe { &*signature };
    let pk = unsafe { &*public_key };
    let msg_str = unsafe { CStr::from_ptr(message) };
    
    let message_obj = Message::new(b"", msg_str.to_bytes());
    let is_valid = sig.signature.verify(&message_obj, &pk.public_key);
    
    if is_valid {
        BLS_SUCCESS
    } else {
        BLS_ERROR_SIGNATURE_VERIFICATION_FAILED
    }
}

// ===== 序列化功能 =====

// 序列化公鑰為字串
#[no_mangle]
pub extern "C" fn bls_public_key_to_string(public_key: *const BLSPublicKey) -> *mut c_char {
    if public_key.is_null() {
        return ptr::null_mut();
    }
    
    let pk = unsafe { &*public_key };
    let bytes = pk.public_key.to_bytes();
    let hex_string = hex::encode(bytes);
    match CString::new(hex_string) {
        Ok(c_str) => c_str.into_raw(),
        Err(_) => ptr::null_mut(),
    }
}

// 序列化私鑰為字串
#[no_mangle]
pub extern "C" fn bls_secret_key_to_string(secret_key: *const BLSSecretKey) -> *mut c_char {
    if secret_key.is_null() {
        return ptr::null_mut();
    }
    
    let sk = unsafe { &*secret_key };
    let bytes = sk.secret_key.to_bytes();
    let hex_string = hex::encode(bytes);
    match CString::new(hex_string) {
        Ok(c_str) => c_str.into_raw(),
        Err(_) => ptr::null_mut(),
    }
}

// 序列化簽名為字串
#[no_mangle]
pub extern "C" fn bls_signature_to_string(signature: *const BLSSignature) -> *mut c_char {
    if signature.is_null() {
        return ptr::null_mut();
    }
    
    let sig = unsafe { &*signature };
    let bytes = sig.signature.to_bytes();
    let hex_string = hex::encode(bytes);
    match CString::new(hex_string) {
        Ok(c_str) => c_str.into_raw(),
        Err(_) => ptr::null_mut(),
    }
}

// ===== 記憶體管理 =====

// 釋放公鑰記憶體
#[no_mangle]
pub extern "C" fn bls_free_public_key(public_key: *mut BLSPublicKey) -> c_int {
    if public_key.is_null() {
        return BLS_ERROR_INVALID_INPUT;
    }
    
    unsafe {
        let _ = Box::from_raw(public_key);
    }
    BLS_SUCCESS
}

// 釋放私鑰記憶體
#[no_mangle]
pub extern "C" fn bls_free_secret_key(secret_key: *mut BLSSecretKey) -> c_int {
    if secret_key.is_null() {
        return BLS_ERROR_INVALID_INPUT;
    }
    
    unsafe {
        let _ = Box::from_raw(secret_key);
    }
    BLS_SUCCESS
}

// 釋放簽名記憶體
#[no_mangle]
pub extern "C" fn bls_free_signature(signature: *mut BLSSignature) -> c_int {
    if signature.is_null() {
        return BLS_ERROR_INVALID_INPUT;
    }
    
    unsafe {
        let _ = Box::from_raw(signature);
    }
    BLS_SUCCESS
}

// 釋放字串記憶體
#[no_mangle]
pub extern "C" fn bls_free_string(s: *mut c_char) -> c_int {
    if s.is_null() {
        return BLS_ERROR_INVALID_INPUT;
    }
    
    unsafe {
        let _ = CString::from_raw(s);
    }
    BLS_SUCCESS
}

// ===== 記憶體管理工具函數 =====

// 檢查指針是否為空
#[no_mangle]
pub extern "C" fn bls_is_null(ptr: *const std::ffi::c_void) -> c_int {
    if ptr.is_null() {
        1
    } else {
        0
    }
}

// 獲取記憶體使用統計 (僅在 debug 模式下)
#[cfg(debug_assertions)]
#[no_mangle]
pub extern "C" fn bls_get_memory_stats() -> *mut c_char {
    // 這裡可以添加記憶體統計功能
    let stats = "Memory stats not implemented yet";
    match CString::new(stats) {
        Ok(c_str) => c_str.into_raw(),
        Err(_) => ptr::null_mut(),
    }
}

#[cfg(not(debug_assertions))]
#[no_mangle]
pub extern "C" fn bls_get_memory_stats() -> *mut c_char {
    ptr::null_mut()
} 