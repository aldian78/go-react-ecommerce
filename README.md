🛒 gocommerce — E-Commerce Backend with Go + gRPC

⚡ gocommerce adalah backend e-commerce modern berbasis Golang + gRPC dengan dukungan PostgreSQL dan Docker.
Dirancang untuk performa tinggi, modular, dan siap diintegrasikan dengan frontend (React).

✨ Fitur Utama

🔑 Auth Service (login, register, JWT)

🛍️ Order Service (CRUD order, transaksi dengan DB transaction)

📦 Product Service (list, detail, add, update, delete)

🗄️ Database PostgreSQL dengan migration support

🐳 Docker + docker-compose untuk local & production deploy

📡 gRPC & Protobuf untuk komunikasi cepat & type-safe

🌐 REST Gateway untuk akses dari client HTTP

📦 gocommerc/
 ┣ 📂 cmd/          → entrypoint service (auth, order, product)
 ┣ 📂 internal/     → business logic & repository
 ┣ 📂 proto/        → protobuf definitions
 ┣ 📂 migrations/   → database migrations
 ┣ 📜 docker-compose.yml
 ┣ 📜 Dockerfile
 ┗ 📜 README.md

👉 Ini akan menjalankan:

gRPC service

PostgreSQL

📧 Kontak

📩 Email: yourname@example.com
🐙 GitHub: username

✨ Selamat berkoding & semoga project ini membantu membangun e-commerce modern dengan Go + gRPC 🚀
