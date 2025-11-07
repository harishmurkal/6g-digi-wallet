# 6g-digi-wallet
Simple reference implementation for Digital-Wallet used w.r.t Telco domain

# build steps
go mod init github.com/harishmurkal/6g-digi-wallet
Remove-Item .\bin\wallet-server.exe
cd 'C:\Users\harism\6g-digi-wallet'
go build -o bin/wallet-server.exe ./cmd/wallet-server

# Execution steps
cd 'C:\Users\harism\6g-digi-wallet'
$env:STORE_BACKEND="file"
.\bin\wallet-server.exe

## Testing 

### Run all
cd 'C:\Users\harism\6g-digi-wallet'
Set-ExecutionPolicy Bypass -Scope Process -Force
.\tests\run-tests.ps1

### DID testing (Issuer)
cd 'C:\Users\harism\6g-digi-wallet'

#### Generate a new DID for issuer and subscriber (e.g. telco:airtel)
curl -Method POST -Uri http://localhost:8080/issuer/did/generate `
  -ContentType "application/json" `
  -InFile .\tests\test-did-generate-airtel.json

curl -Method POST -Uri http://localhost:8080/issuer/did/generate `
  -ContentType "application/json" `
  -InFile .\tests\test-did-generate-harism.json

#### Resolve DID
curl http://localhost:8080/issuer/did/did:telco:airtel
curl http://localhost:8080/issuer/did/did:telco:harism

#### Create VC via issuer (POST with a VCRequest JSON, returns signed VC)
curl -Method POST -Uri http://localhost:8080/issuer/vc/create `
  -ContentType "application/json" `
  -InFile .\tests\test-vc-request-IDAndLoc.json > .\tests\tmp_signed_vc.json

### VC testing (Issuer + Wallet)

#### Store DID in wallet
curl -Method POST -Uri http://localhost:8080/wallet/did/store `
  -ContentType "application/json" `
  -InFile .\tests\test-did-store.json

#### Store VC in wallet
curl -Method POST -Uri http://localhost:8080/wallet/vc/store `
  -ContentType "application/json" `
  -InFile .\tests\test-vc-store.json

#### Get specific DID from wallet
curl http://localhost:8080/wallet/did/did:telco:airtel

#### Retrieve a VC from wallet
curl http://localhost:8080/wallet/vc/vc-harism-sim001

#### List stored DIDs in wallet
curl http://localhost:8080/wallet/did/list

#### List all VCs in wallet
curl http://localhost:8080/wallet/vc/list

#### Build a VP (Verifiable Presentation) at Wallet
curl -Method POST -Uri http://localhost:8080/wallet/vp/build `
  -ContentType "application/json" `
  -InFile .\tests\test-vp-build.json

#### Verify a VC (wallet end verification of VC)
curl -Method POST -Uri http://localhost:8080/wallet/verify `
  -ContentType "application/json" `
  -InFile .\tests\test-vc-verify.json

### VP testing (Wallet + Verifier)
#### Verify the VP at verifier
curl -Method POST -Uri http://localhost:8080/verifier/vp/verify `
  -ContentType "application/json" `
  -InFile .\tests\test-vp-verify.json

# Version information
PS C:\Users\harism\6g-digi-wallet> go version
go version go1.25.0 windows/amd64
PS C:\Users\harism\6g-digi-wallet>

# Debugging information
PS C:\Users\harism\6g-digi-wallet> go mod init github.com/harishmurkal/6g-digi-wallet
go: creating new go.mod: module github.com/harishmurkal/6g-digi-wallet
go: to add module requirements and sums:
        go mod tidy
PS C:\Users\harism\6g-digi-wallet> go mod tidy
PS C:\Users\harism\6g-digi-wallet>

# PlantUML generation
cd 'C:\Users\harism\6g-digi-wallet\docs\PlantUMLs'
java -jar 'C:\Users\harism\OneDrive - Nokia\D\UTILITIES\plantuml-1.2025.7.jar' did-vc-vp-flow.puml