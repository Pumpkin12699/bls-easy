.PHONY: build test clean

# 編譯 Rust 庫
build:
	cargo build --release
	cbindgen --config cbindgen.toml --crate w3f-bls-easy --output src/bls_easy_ffi.h

# 測試 Go 程式
test: build
	go run main.go

# 清理
clean:
	cargo clean
	rm -f src/bls_easy_ffi.h

# 安裝 cbindgen (如果還沒安裝)
install-cbindgen:
	cargo install --force cbindgen

# 完整設置
setup: install-cbindgen build test 