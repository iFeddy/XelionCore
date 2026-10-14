# XelionCore ⚡

**XelionCore** es una reimplementación moderna, ultra liviana y multiplataforma (Linux, macOS, Windows) del servidor de **Fiesta Online** desarrollada en **Go (Golang)**.

Diseñado desde cero para resolver el problema de consumo masivo de memoria de los ejecutables legados de Windows (que requerían ~24 GB de RAM), funcionando en un solo proceso compilado estático con un consumo menor a **150 MB de RAM**.

---

## 🚀 Características

* **Rendimiento Máximo y Cero Dependencias:** Compila a un binario estático único sin necesidad de runtimes externos ni máquinas virtuales.
* **Criptografía Nativa:** Implementación exacta de la tabla XOR rolling de 499 bytes (`NetCrypto`).
* **Protocolo Verificado:** Opcodes y estructuras de paquetes alineadas con la versión **NA2016** (TeamNG).
* **Multiplataforma Real:** Corre nativamente en macOS (Apple Silicon M1/M2/M3/M4 e Intel), servidores Linux (Ubuntu/Debian) y Windows.
* **Persistencia Flexible:** Incluye un repositorio en memoria (`MemoryAccountRepository`) para arrancar de inmediato sin configurar bases de datos, con soporte para MySQL/MariaDB.

---

## 📁 Estructura del Proyecto

```text
XelionCore/
├── cmd/
│   └── login/             # Punto de entrada del Login Server
├── configs/
│   └── config.json        # Configuración de puertos y mundos
├── pkg/
│   ├── crypto/            # NetCrypto (tabla XOR de 499 bytes)
│   ├── database/          # Modelos de cuenta y repositorios (Memory / MySQL)
│   ├── network/           # TCP Server, sesiones de cliente y framing
│   ├── protocol/          # PacketReader, PacketWriter y Opcodes
│   └── server/
│       └── login/         # Lógica y Handlers del Login Server
├── go.mod
└── README.md
```

---

## 🛠️ Cómo compilar y ejecutar

### Requisitos
* [Go 1.22+](https://go.dev/dl/) instalado en tu sistema.

### Compilar y correr el Login Server:
```bash
# Descargar dependencias
go mod tidy

# Ejecutar directamente
go run ./cmd/login

# O compilar binario
go build -o bin/login ./cmd/login
./bin/login
```

### Cuenta por defecto:
El servidor viene preconfigurado con una cuenta de pruebas en memoria:
* **Usuario:** `admin`
* **Contraseña:** `admin`
* *Cualquier usuario nuevo que intente loguearse se auto-crea automáticamente.*

---

## 📜 Licencia
Proyecto distribuido bajo la licencia MIT.
