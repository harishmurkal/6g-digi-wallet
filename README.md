# 6g-digi-wallet

Simple reference implementation for Digital-Wallet used w.r.t Telco domain

## 💡 Project Overview

This repository contains the foundational implementation of a **Digital Identity Wallet** and associated **Issuer/Verifier Services** compliant with W3C Verifiable Credentials (VC) and Decentralized Identifiers (DID) standards.

The primary goal is to demonstrate a secure, privacy-preserving identity layer suitable for future 6G-era telecom services, where immutable subscriber identity and selective data disclosure are paramount. The implementation is written in Go and structured around the three core actors in a decentralized identity ecosystem: the **Issuer**, the **Holder (Wallet)**, and the **Verifier**.

---

## 🏗️ Architecture and Components

The system is organized into three logical services and four core data structures, all managed through a persistent (in-memory) store.

### Services

| Service | Role | Core Functionality |
| :--- | :--- | :--- |
| `Issuer Service` | **Telco Operator** | Generates DID Documents, creates and signs Verifiable Credentials (VCs). |
| `Wallet Service` | **Subscriber (Holder)** | Securely stores VCs, constructs and signs Verifiable Presentations (VPs) based on verifier requests. |
| `Verifier Service`| **Relying Party** | Resolves DIDs, verifies the cryptographic integrity of VPs and embedded VCs, and enforces policy checks. |

### Core Data Structures (`internal/models`)

| Structure | Description |
| :--- | :--- |
| `DIDDocument` | The public document describing the cryptographic key material associated with an identity. |
| `VerifiableCredential` (VC) | A tamper-proof claim signed by the Issuer (e.g., "This DID owns SIM X"). |
| `VerifiablePresentation` (VP)| A container created and signed by the Holder (Wallet) to selectively share VCs with a Verifier. |
| `CryptoService` | An abstract interface centralizing all key generation, signing, and verification (Ed25519/JWS). |

---

## 🔄 Digital Identity Flow

The entire lifecycle—from key generation to authentication—follows a standardized pattern.

The diagram below illustrates the sequence, highlighting the cryptographic steps taken by the Issuer (signing VCs) and the Holder (signing VPs).

![Sequence Diagram showing the flow between Issuer, Holder, and Verifier](docs/PlantUMLs/flow-crypto6g.png)

---

## 🛠️ Setup and Build

### Prerequisites
* Go 1.21+ (Tested with Go 1.25.0)
* The `uuid` package (`github.com/google/uuid`)

### Build Steps

```bash
# Initialize Go module
go mod init github.com/harishmurkal/6g-digi-wallet
go mod tidy

# Build the wallet server executable
Remove-Item .\bin\wallet-server.exe -ErrorAction SilentlyContinue
cd 'C:\Users\harism\6g-digi-wallet'
go build -o bin/wallet-server.exe ./cmd/wallet-server
```

### Execution Steps

```bash
cd 'C:\Users\harism\6g-digi-wallet'

# Set storage backend (options: "memory" or "file")
$env:STORE_BACKEND="file"

# Run the wallet server
.\bin\wallet-server.exe
```

---

## 🧪 Testing

### Unit Testing

```bash
cd 'C:\Users\harism\6g-digi-wallet'

# Test specific modules
go test ./internal/service/crypto6g -v
go test ./internal/storage -v

# Run all tests
go test ./... -v
```

### Functional Testing

```bash
cd 'C:\Users\harism\6g-digi-wallet'
Set-ExecutionPolicy Bypass -Scope Process -Force
.\tests\run-tests.ps1
```

---

## 📡 API Usage Examples

### DID Testing (Issuer)

#### Generate DIDs for Issuer and Subscriber

```powershell
cd 'C:\Users\harism\6g-digi-wallet'

# Generate DID for issuer (e.g., telco:airtel)
curl -Method POST -Uri http://localhost:8080/issuer/did/generate `
  -ContentType "application/json" `
  -InFile .\tests\test-did-generate-airtel.json

# Generate DID for subscriber (e.g., telco:harism)
curl -Method POST -Uri http://localhost:8080/issuer/did/generate `
  -ContentType "application/json" `
  -InFile .\tests\test-did-generate-harism.json
```

#### Resolve DIDs

```powershell
curl http://localhost:8080/issuer/did/did:telco:airtel
curl http://localhost:8080/issuer/did/did:telco:harism
```

#### Create Verifiable Credential

```powershell
# Create and sign a VC via issuer
curl -Method POST -Uri http://localhost:8080/issuer/vc/create `
  -ContentType "application/json" `
  -InFile .\tests\test-vc-request-IDAndLoc.json > .\tests\tmp_signed_vc.json
```

### VC Testing (Issuer + Wallet)

#### Store DID and VC in Wallet

```powershell
# Store DID in wallet
curl -Method POST -Uri http://localhost:8080/wallet/did/store `
  -ContentType "application/json" `
  -InFile .\tests\test-did-store.json

# Store VC in wallet
curl -Method POST -Uri http://localhost:8080/wallet/vc/store `
  -ContentType "application/json" `
  -InFile .\tests\test-vc-store.json
```

#### Retrieve from Wallet

```powershell
# Get specific DID
curl http://localhost:8080/wallet/did/did:telco:airtel

# Retrieve specific VC
curl http://localhost:8080/wallet/vc/vc-harism-sim001

# List all DIDs
curl http://localhost:8080/wallet/did/list

# List all VCs
curl http://localhost:8080/wallet/vc/list
```

#### Build Verifiable Presentation

```powershell
# Build VP at wallet (selective disclosure)
curl -Method POST -Uri http://localhost:8080/wallet/vp/build `
  -ContentType "application/json" `
  -InFile .\tests\test-vp-build.json
```

#### Verify Credential at Wallet

```powershell
# Wallet-side VC verification
curl -Method POST -Uri http://localhost:8080/wallet/verify `
  -ContentType "application/json" `
  -InFile .\tests\test-vc-verify.json
```

### VP Testing (Wallet + Verifier)

#### Verify Presentation at Verifier

```powershell
# Verify VP at verifier endpoint
curl -Method POST -Uri http://localhost:8080/verifier/vp/verify `
  -ContentType "application/json" `
  -InFile .\tests\test-vp-verify.json
```

---

