Here is a highly professional, comprehensive `README.md` perfectly tailored to the exact structure of your repository. It highlights the new modular API design, the specific Go collectors, and the enterprise-grade security features.

You can copy and paste this directly into the root `README.md` file of your GitHub repository.

---

# 🛡️ Vader-Pro-VDR (Vulnerability Detection & Response)

**Vader-Pro-VDR** is a lightweight, distributed endpoint telemetry and asset management pipeline designed for secure, asynchronous data collection across remote networks. Built with a Zero-Trust mindset, the platform decouples localized system profiling from centralized ingestion, creating a resilient "spool-and-ship" architecture that guarantees zero data loss—even in highly restrictive or unstable network environments.

---

## 🏗️ Repository Structure

The project is divided into two primary, decoupled components:

* **`/vmagent` (The Edge Endpoint):** A highly concurrent, zero-dependency Go binary deployed to target machines. It maps host telemetry, encrypts the payload locally, and manages reliable outbound network transmission.
* **`/Python_Server-API` (The Central Core):** A modular, asynchronous Python FastAPI gateway. It authenticates inbound agents, decrypts AES-GCM payloads in memory, and stores structured telemetry in a SQLite database for dashboarding and analysis.

---

## ✨ Key Architectural Features

### 1. Multi-OS Telemetry Collection

The Go agent utilizes `gopsutil` and native OS commands to seamlessly profile endpoints across environments. It collects Hostname, OS details, CPU architecture, Memory utilization, Network Interfaces (stripping loopbacks and classifying types), Disk mount usage, and full Software Inventories (via Windows Registry, `dpkg`, `rpm`, or macOS `system_profiler`).

### 2. Secure Store-and-Forward (Spooling)

The agent operates on a dual-routine architecture to prevent "telemetry lag." A background scanner profiles the host (every 7 days) and instantly encrypts the payload locally to disk (`inventory.enc`). A detached shipper loop (checking every 10 seconds) monitors this spool and transmits it when the network is available, clearing the disk upon successful server acknowledgment.

### 3. End-to-End Payload Encryption

Transmissions are mathematically locked using **AES-256-GCM** before ever touching the transport layer. This ensures that even if SSL/TLS is stripped at perimeter firewalls, the system profile data remains entirely unreadable to network sniffers.

### 4. Modular API Gateway

The central server is structured for enterprise scalability. `main.py` acts as a clean entry point with CORS enabled, delegating traffic to specific routers:

* `routers/ingestion.py`: Handles secure `POST` payloads from edge agents.
* `routers/dashboard.py`: Exposes `GET` endpoints to serve decrypted telemetry to future frontend UI panels.

---

## 🛠️ Technology Stack

**Edge Agent (`vmagent`)**

* **Language:** Go 1.20+
* **Concurrency:** Native Goroutines (Decoupled Scanner & Shipper loops)
* **Cryptography:** `crypto/aes` & `crypto/cipher` (AES-256-GCM)
* **System Parsing:** `shirou/gopsutil`

**Central Management Server (`Python_Server-API`)**

* **Framework:** Python / FastAPI (Asynchronous ASGI)
* **Web Server:** Uvicorn
* **Cryptography:** `pycryptodome`
* **Database:** SQLite3 (WAL mode enabled, utilizing native JSON data structures)

---

## 🚀 Getting Started

### Prerequisites

* **Server:** Python 3.8+, `pip`, `venv`
* **Agent:** Go (for compiling the binary)

### 1. Start the Central Receiver

Navigate to the server directory, set up the Python environment, and launch the API:

```bash
cd Python_Server-API

# Create and activate virtual environment
python3 -m venv venv
source venv/bin/activate  # On Windows use: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Launch the ingestion API (runs on port 8000)
uvicorn main:app --host 0.0.0.0 --port 8000

```

### 2. Deploy the Edge Agent

Before compiling, ensure the `serverURL` in `vmagent/internal/transport/http.go` (Line 14) points to your active central server IP.

On the target endpoint, navigate to the agent directory, compile, and run:

```bash
cd vmagent

# Sync dependencies
go mod tidy

# Build the executable
go build -o vmagent cmd/main.go

# Run the agent
./vmagent

```

*(The agent will immediately profile the system, encrypt the data, and securely ship it to the waiting API).*

---

## 🗺️ Roadmap

* [x] **Phase 1:** Multi-OS endpoint hardware and software baseline collection.
* [x] **Phase 2:** Secure asynchronous transport, local spooling, and AES-GCM encryption pipeline.
* [x] **Phase 2.5:** Modular API restructure (separating ingestion and dashboard logic).
* [ ] **Phase 3:** Frontend UI Dashboard (Single Pane of Glass).
* [ ] **Phase 4:** Threat Intelligence Correlation (Integrating OSV and NVD for real-time CVE mapping).
