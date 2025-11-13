# Build Stage - "เตาหลอม" สำหรับ compile code
# เราใช้ image 'golang' ที่มีเครื่องมือ Go ทั้งหมด (ตามที่อาจารย์ระบุ)
FROM golang:1.24.5 AS builder

# ตั้งค่าโฟลเดอร์ทำงาน
WORKDIR /app

# 1. คัดลอก "ใบรายการ" (go.mod, go.sum) ก่อน
# (นี่คือเทคนิค Layer Caching ที่อาจารย์สอน)
COPY go.mod go.sum ./
# 2. โหลด dependencies (ถ้า go.mod ไม่เปลี่ยน, Docker จะใช้ cache)
RUN go mod download

# 3. คัดลอก "โค้ดทั้งหมด" ของเรา
COPY . .

# 4. "หลอม" (Build) โค้ดของเราให้เป็นไฟล์ binary
# (นี่คือคำสั่ง 'static binary' ที่อาจารย์สอนเป๊ะๆ)
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main ./cmd/api/main.go

# ---

# Run Stage - "กล่อง" ที่จะเอาไปใช้งานจริง
# เราใช้ 'alpine' ที่เล็กและปลอดภัย (ตามที่อาจารย์สอน)
FROM alpine:latest  

# (จำเป็นสำหรับ HTTPS/SSL, ตามที่อาจารย์สอน)
RUN apk --no-cache add ca-certificates curl

# ตั้งค่าโฟลเดอร์ทำงาน
WORKDIR /app/

# 5. คัดลอกจาก "เตาหลอม" (builder) ... เอามาแค่ "ดาบ" (ไฟล์ main)
COPY --from=builder /app/main .

# เปิด port ที่ Go API ของเรา (main.go) ใช้อยู่
EXPOSE 8080

# 6. คำสั่ง "เปิด" กล่อง
# รันไฟล์ binary ที่เรา build มา
ENTRYPOINT ["./main"]