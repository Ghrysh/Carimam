# 🛒 CariMam - Microservices Backend API

CariMam adalah sebuah platform pemesanan makanan (*Food Delivery*) yang dibangun menggunakan arsitektur **Microservices**. Proyek ini mendemonstrasikan komunikasi antar-layanan, manajemen otentikasi terpusat, dan penggunaan API Gateway sebagai *Single Entry Point*.

## 🏗️ Arsitektur Sistem
Sistem ini terdiri dari 4 layanan utama yang saling berkolaborasi:

1. **API Gateway (Port 8000):** *Reverse Proxy* untuk merutekan *request* dari *client* ke *service* yang tepat.
2. **Identity Service (Port 8080):** Menangani Autentikasi (JWT), Otorisasi (Role-based), dan Manajemen Saldo/Dompet (*Balance*).
3. **Product Service (Port 8081):** Menangani katalog menu makanan, stok, dan *upload* gambar.
4. **Order Service (Port 8082):** Menangani *checkout* pesanan, kalkulasi harga (komunikasi HTTP ke Product), dan pemotongan saldo (komunikasi HTTP ke Identity).

## 🚀 Teknologi yang Digunakan
* **Bahasa:** Golang (Go 1.25)
* **Framework:** Gin Web Framework
* **Database:** PostgreSQL (3 Database terpisah untuk masing-masing *service*)
* **ORM:** GORM
* **Infrastruktur:** Docker & Docker Compose
* **Keamanan:** JWT (JSON Web Token) & Bcrypt

## ⚙️ Cara Menjalankan Aplikasi

**1. Jalankan Database (Docker Compose)**
```bash
docker-compose up -d