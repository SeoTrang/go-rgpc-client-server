Dưới đây là hướng dẫn chi tiết cách sử dụng project gRPC với Go, được biên soạn thành file `README.md`:

---

# gRPC Server và Client trong Go

Đây là ví dụ đơn giản về cách xây dựng một gRPC Server và Client trong Go. Server cung cấp một service tên là `Greeter` với một method `SayHello`, và Client gọi method này để nhận phản hồi.

## Mục Lục
- [Giới Thiệu](#giới-thiệu)
- [Cài Đặt Công Cụ](#cài-đặt-công-cụ)
- [Cấu Trúc Thư Mục](#cấu-trúc-thư-mục)
- [Tạo File `greeter.proto`](#tạo-file-greeterproto)
- [Generate Code từ Proto](#generate-code-từ-proto)
- [Tạo gRPC Server](#tạo-grpc-server)
- [Tạo gRPC Client](#tạo-grpc-client)
- [Chạy Project](#chạy-project)

---

## Giới Thiệu
Mục tiêu của dự án này là xây dựng một gRPC server và client đơn giản:

- **Server**: Tạo một service tên `Greeter` với một method `SayHello`, nhận tên người dùng và trả về một message chào họ.
- **Client**: Gửi tên người dùng tới server và nhận phản hồi với message chào từ server.

---

## Cài Đặt Công Cụ

### 1. Cài `protoc` (Protocol Buffers Compiler)
Trước tiên, bạn cần cài đặt `protoc` để biên dịch các file `.proto`.

- **Mac**: 
  ```bash
  brew install protobuf
  ```
- **Ubuntu**: 
  ```bash
  sudo apt install -y protobuf-compiler
  ```
- **Windows**: 
  Cài đặt qua [Chocolatey](https://chocolatey.org) bằng lệnh:
  ```bash
  choco install protoc
  ```

### 2. Cài plugin Go cho `protoc`
Chạy lệnh dưới đây để cài plugin Go cho `protoc`:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

**Lưu ý**: Đảm bảo thư mục `$GOPATH/bin` nằm trong `$PATH`.

---

## Cấu Trúc Thư Mục

```
grpc-example/
├── go.mod
├── greeter.proto
├── server/
│   └── main.go
├── client/
│   └── main.go
```

---

## Tạo File `greeter.proto`

File `greeter.proto` định nghĩa giao diện của gRPC service, bao gồm method `SayHello`.

```proto
syntax = "proto3";

package greeter;

option go_package = "./greeterpb";

service Greeter {
  rpc SayHello (HelloRequest) returns (HelloResponse);
}

message HelloRequest {
  string name = 1;
}

message HelloResponse {
  string message = 1;
}
```

**Giải thích**:
- **service Greeter**: Đây là dịch vụ gRPC với method `SayHello`.
- **message HelloRequest**: Message gửi từ client, chứa tên người dùng.
- **message HelloResponse**: Message trả về từ server, chứa lời chào.

---

## Generate Code từ Proto

Sau khi tạo file `greeter.proto`, bạn cần sử dụng `protoc` để sinh mã Go từ file `.proto`:

```bash
protoc --go_out=. --go-grpc_out=. greeter.proto
```

Lệnh trên sẽ tạo ra thư mục `greeterpb/` chứa các file Go tự động sinh ra từ file `greeter.proto`.

---

## Tạo gRPC Server

File `server/main.go` triển khai server gRPC:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"grpc-example/greeterpb"

	"google.golang.org/grpc"
)

type server struct {
	greeterpb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *greeterpb.HelloRequest) (*greeterpb.HelloResponse, error) {
	name := req.GetName()
	res := &greeterpb.HelloResponse{
		Message: "Hello " + name,
	}
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greeterpb.RegisterGreeterServer(s, &server{})

	fmt.Println("🚀 Server is running at :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```

**Giải thích**:
- Server khởi tạo lắng nghe ở cổng `50051`.
- Khi nhận request từ client, method `SayHello` sẽ trả về một message chứa lời chào.

---

## Tạo gRPC Client

File `client/main.go` thực hiện kết nối tới server và gọi method `SayHello`:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"grpc-example/greeterpb"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("❌ could not connect: %v", err)
	}
	defer conn.Close()

	client := greeterpb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.SayHello(ctx, &greeterpb.HelloRequest{Name: "Go Developer"})
	if err != nil {
		log.Fatalf("❌ could not greet: %v", err)
	}
	fmt.Println("👋 Response:", res.GetMessage())
}
```

**Giải thích**:
- Client kết nối tới server gRPC tại `localhost:50051`.
- Client gọi method `SayHello` và in ra phản hồi từ server.

---

## Chạy Project

Để chạy project, bạn cần mở hai terminal:

### Terminal 1: Chạy Server
```bash
go run server/main.go
```

### Terminal 2: Chạy Client
```bash
go run client/main.go
```

Kết quả:

```
👋 Response: Hello Go Developer
```

---

## Giải Thích Thêm

### Một số đơn vị thời gian phổ biến:
- **ms** (Millisecond): 1 ms = 0.001 giây.
- **µs** (Microsecond): 1 µs = 0.000001 giây.
- **ns** (Nanosecond): 1 ns = 0.000000001 giây.

Thời gian phản hồi của gRPC rất nhanh, chỉ khoảng **500µs đến 700µs**, cho thấy gRPC là một giao thức tối ưu cho việc xử lý các yêu cầu nội bộ.

---

### So sánh gRPC trong Go và Node.js

- **Node.js** sử dụng **Dynamic Code Generation** và **Reflection** để làm việc với các file `.proto` trực tiếp tại runtime. Điều này giúp giảm bước build nhưng có thể làm chậm hiệu suất và không kiểm tra lỗi mạnh mẽ.
- **Go** yêu cầu **build** các file `.proto` thành mã Go trước khi biên dịch, giúp tối ưu hiệu suất và kiểm tra lỗi ngay khi biên dịch.

---