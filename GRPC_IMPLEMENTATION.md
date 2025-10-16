# Auth Service dengan gRPC - Implementasi Lengkap

## 🎉 Implementasi Berhasil Diselesaikan!

Anda sekarang memiliki auth service yang mendukung **dual protocol**:
- ✅ **REST API** (port 9000)
- ✅ **gRPC API** (port 9001)

## 🚀 Fitur yang Diimplementasikan

### 1. **gRPC Services**
- **Register** - Registrasi user baru
- **Login** - Autentikasi user
- **RefreshToken** - Refresh access token
- **Logout** - Logout user
- **ChangePassword** - Ubah password
- **CheckAccess** - Cek permission akses
- **GetUserProfile** - Ambil profil user
- **GetUserFunctions** - Ambil fungsi/permission user

### 2. **Autentikasi JWT**
- Support JWT untuk gRPC melalui metadata
- Kompatibel dengan sistem autentikasi REST API yang sudah ada
- Reuse logic authentication yang sudah ada

### 3. **Protocol Buffer Definition**
- File `proto/auth.proto` dengan semua message dan service definitions
- Generated Go code di `grpc/pb/`

### 4. **gRPC Server Implementation**
- Server implementation di `grpc/server.go`
- Adapter untuk reuse usecase layer yang sudah ada
- Error handling yang proper dengan gRPC status codes

## 📁 Struktur File yang Ditambahkan

```
auth-service/
├── proto/
│   └── auth.proto              # Protocol buffer definition
├── grpc/
│   ├── pb/
│   │   ├── auth.pb.go         # Generated protobuf messages
│   │   └── auth_grpc.pb.go    # Generated gRPC service
│   ├── server.go              # gRPC server implementation
│   └── README.md              # gRPC documentation
├── examples/
│   ├── grpc_client.go         # gRPC client example
│   └── run_demo.sh            # Demo script
├── Makefile                   # Build automation
└── main.go                    # Updated untuk dual server
```

## 🔧 Cara Menjalankan

### 1. **Jalankan Server (REST + gRPC)**
```bash
make run
# atau
go run main.go
```

### 2. **Test gRPC Client**
```bash
make test-grpc
# atau
go run examples/grpc_client.go
```

### 3. **Demo Lengkap**
```bash
make demo
# atau
./examples/run_demo.sh
```

## 🌐 Endpoint yang Tersedia

### REST API (Port 9000)
- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/refresh`
- `POST /api/auth/logout`
- `GET /api/auth/me`
- dan lainnya...

### gRPC API (Port 9001)
- `auth.AuthService/Register`
- `auth.AuthService/Login`
- `auth.AuthService/RefreshToken`
- `auth.AuthService/Logout`
- `auth.AuthService/GetUserProfile`
- `auth.AuthService/CheckAccess`
- `auth.AuthService/GetUserFunctions`

## 🧪 Testing Tools

### 1. **Go Client**
```bash
go run examples/grpc_client.go
```

### 2. **grpcurl** (jika terinstall)
```bash
# List services
grpcurl -plaintext localhost:9001 list

# Test register
grpcurl -plaintext -d '{
  "full_name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "gender": 1
}' localhost:9001 auth.AuthService/Register
```

### 3. **GUI Tools**
- BloomRPC
- Postman (dengan gRPC support)
- Insomnia

## 🐳 Docker Support

### Build dan Run dengan Docker
```bash
make docker-build
make docker-run
```

### Docker Compose
```bash
docker-compose up
```

## 📚 Dependencies yang Ditambahkan

- `google.golang.org/grpc` - gRPC framework
- `google.golang.org/protobuf` - Protocol buffers

## 🎯 Keunggulan Implementasi

1. **Dual Protocol**: Klien bisa pilih REST atau gRPC sesuai kebutuhan
2. **Reuse Business Logic**: Tidak ada duplikasi kode, semua menggunakan usecase layer yang sama
3. **JWT Authentication**: Support autentikasi untuk kedua protocol
4. **Error Handling**: Proper error handling untuk gRPC dengan status codes
5. **Documentation**: Dokumentasi lengkap dan contoh penggunaan
6. **Testing**: Client example dan demo script
7. **Development Tools**: Makefile untuk automation

## 🔄 Regenerate Proto Files

Jika Anda mengubah `proto/auth.proto`:
```bash
make proto
```

## 📖 Dokumentasi Lebih Lanjut

- [gRPC README](grpc/README.md) - Dokumentasi detail gRPC implementation
- [Main README](README.md) - Dokumentasi utama project

---

🎉 **Auth service Anda sekarang support gRPC!** 
Kedua server (REST dan gRPC) berjalan bersamaan dan siap digunakan.