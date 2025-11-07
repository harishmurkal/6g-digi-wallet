<#
.SYNOPSIS
    6G Digi Wallet End-to-End API Test Runner (PowerShell native)

.DESCRIPTION
    Executes API tests for Issuer, Wallet, and Verifier using Invoke-RestMethod.
    Captures DID Documents automatically for Wallet storage tests.
#>

param(
    [switch]$OnlyIssuer,
    [switch]$OnlyWallet,
    [switch]$OnlyVerifier,
    [switch]$Help
)

if ($Help) {
    Get-Help $MyInvocation.MyCommand.Path -Detailed
    exit 0
}

$ErrorActionPreference = "Stop"

# ---------------------------------------------------------------------
# Setup paths
# ---------------------------------------------------------------------
$BaseDir = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
$OutputDir = Join-Path $BaseDir "output"
$LogFile   = Join-Path $OutputDir "run-log.txt"

if (!(Test-Path $OutputDir)) {
    New-Item -ItemType Directory -Force -Path $OutputDir | Out-Null
}

"=== 6G Digi Wallet Test Run ($(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')) ===`n" |
    Out-File -FilePath $LogFile -Encoding utf8 -Force

# ---------------------------------------------------------------------
# Helper: Log
# ---------------------------------------------------------------------
function Log-Write {
    param($Text)
    $Text | Tee-Object -FilePath $LogFile -Append
}

# ---------------------------------------------------------------------
# Helper: Run-Step using Invoke-RestMethod
# ---------------------------------------------------------------------
function Run-Step {
    param (
        [string]$Description,
        [string]$Method,
        [string]$Url,
        [string]$InFile = "",
        [string]$CaptureTo = ""
    )

    $timestamp = Get-Date -Format "HH:mm:ss"
    Write-Host "[$timestamp] $Description ..." -ForegroundColor Cyan
    Add-Content -Path $LogFile -Value "`n[$timestamp] START: $Description"

    try {
        $headers = @{ "Content-Type" = "application/json" }

        if ($Method -eq "GET") {
            $resp = Invoke-RestMethod -Method Get -Uri $Url -Headers $headers -ErrorAction Stop
            $outString = if ($null -ne $resp) { $resp | ConvertTo-Json -Depth 10 } else { "" }
        }
        else {
            $body = ""
            if ($InFile -ne "") {
                $fullPath = Join-Path $BaseDir $InFile
                if (-not (Test-Path $fullPath)) {
                    throw "Input file not found: $fullPath"
                }
                $body = Get-Content -Raw -Path $fullPath -ErrorAction Stop
            }
            $resp = Invoke-RestMethod -Method $Method -Uri $Url -Headers $headers -Body $body -ErrorAction Stop
            $outString = if ($null -ne $resp) { $resp | ConvertTo-Json -Depth 10 } else { "" }
        }

        # Log output
        if ($outString -ne "") {
            Add-Content -Path $LogFile -Value $outString
        }
        else {
            Add-Content -Path $LogFile -Value "(no JSON body returned)"
        }

        # Capture JSON to file if requested
        if ($CaptureTo -ne "" -and $outString.Trim().StartsWith("{")) {
            $capturePath = Join-Path $BaseDir $CaptureTo
            $outString | Out-File -FilePath $capturePath -Encoding utf8 -Force
            Write-Host "Captured JSON to $capturePath" -ForegroundColor Yellow
            Add-Content -Path $LogFile -Value "Captured JSON -> $capturePath"
        }

        Write-Host "Success: $Description" -ForegroundColor Green
        Add-Content -Path $LogFile -Value "[$timestamp] Success: $Description"
    }
    catch {
        Write-Host "Error during: $Description" -ForegroundColor Red
        Write-Host $_.Exception.Message -ForegroundColor DarkRed
        Add-Content -Path $LogFile -Value "[$timestamp] Error: $Description - $($_.Exception.Message)"
        throw
    }
    finally {
        # just to satisfy PowerShell's try/catch/finally syntax completeness
    }
}

# ---------------------------------------------------------------------
# Test definitions
# ---------------------------------------------------------------------
$IssuerTests = @(
    @{ Desc = "Generate DID (Airtel)"; Method = "POST"; Url = "http://localhost:8080/issuer/did/generate"; InFile = ".\tests\test-did-generate-airtel.json"; CaptureTo = "" },
    @{ Desc = "Generate DID (Harism)"; Method = "POST"; Url = "http://localhost:8080/issuer/did/generate"; InFile = ".\tests\test-did-generate-harism.json"; CaptureTo = "" },
    @{ Desc = "Resolve DID (Airtel)"; Method = "GET"; Url = "http://localhost:8080/issuer/did/did:telco:airtel"; InFile = ""; CaptureTo = ".\tests\output\test-did-store-airtel.json" },
    @{ Desc = "Resolve DID (Harism)"; Method = "GET"; Url = "http://localhost:8080/issuer/did/did:telco:harism"; InFile = ""; CaptureTo = ".\tests\output\test-did-store-harism.json" },
    @{ Desc = "Create VC via Issuer"; Method = "POST"; Url = "http://localhost:8080/issuer/vc/create"; InFile = ".\tests\test-vc-request-IDAndLoc.json"; CaptureTo = ".\tests\output\test-vc-signed.json" }
)

$WalletTests = @(
    @{ Desc = "Store DID (Airtel) in Wallet"; Method = "POST"; Url = "http://localhost:8080/wallet/did/store"; InFile = ".\tests\output\test-did-store-airtel.json"; CaptureTo = "" },
    @{ Desc = "Store DID (Harism) in Wallet"; Method = "POST"; Url = "http://localhost:8080/wallet/did/store"; InFile = ".\tests\output\test-did-store-harism.json"; CaptureTo = "" },
    @{ Desc = "Verify VC at Wallet"; Method = "POST"; Url = "http://localhost:8080/wallet/vc/verify"; InFile = ".\tests\output\test-vc-signed.json"; CaptureTo = "" },
    @{ Desc = "Store VC in Wallet"; Method = "POST"; Url = "http://localhost:8080/wallet/vc/store"; InFile = ".\tests\output\test-vc-signed.json"; CaptureTo = "" },
    @{ Desc = "Get DID from Wallet"; Method = "GET"; Url = "http://localhost:8080/wallet/did/did:telco:airtel"; InFile = ""; CaptureTo = "" },
    @{ Desc = "Get VC from Wallet"; Method = "GET"; Url = "http://localhost:8080/wallet/vc/vc:did:telco:harism:uuid8"; InFile = ""; CaptureTo = "" },
    @{ Desc = "List all Wallet DIDs"; Method = "GET"; Url = "http://localhost:8080/wallet/did/list"; InFile = ""; CaptureTo = "" },
    @{ Desc = "List all Wallet VCs"; Method = "GET"; Url = "http://localhost:8080/wallet/vc/list"; InFile = ""; CaptureTo = "" },
    @{ Desc = "Build VP in Wallet"; Method = "POST"; Url = "http://localhost:8080/wallet/vp/build"; InFile = ".\tests\test-vp-build.json"; CaptureTo = ".\tests\output\test-vp-signed.json" }
    
)

$VerifierTests = @(
    @{ Desc = "Verify VP at Verifier"; Method = "POST"; Url = "http://localhost:8080/verifier/vp/verify"; InFile = ".\tests\output\test-vp-signed.json"; CaptureTo = "" }
)

# ---------------------------------------------------------------------
# Utility: test runner dispatcher
# ---------------------------------------------------------------------
function Run-TestList {
    param([array]$List)
    foreach ($t in $List) {
        Write-Host ("-" * 80) -ForegroundColor DarkGray
        Run-Step -Description $t.Desc -Method $t.Method -Url $t.Url -InFile $t.InFile -CaptureTo $t.CaptureTo
    }
}

# ---------------------------------------------------------------------
# Run Tests based on parameters
# ---------------------------------------------------------------------
try {
    if ($OnlyIssuer) {
        Write-Host "Running only ISSUER tests..." -ForegroundColor Yellow
        Run-TestList -List $IssuerTests
    }
    elseif ($OnlyWallet) {
        Write-Host "Running only WALLET tests..." -ForegroundColor Yellow
        Run-TestList -List $WalletTests
    }
    elseif ($OnlyVerifier) {
        Write-Host "Running only VERIFIER tests..." -ForegroundColor Yellow
        Run-TestList -List $VerifierTests
    }
    else {
        Write-Host "Running ALL test sets (Issuer + Wallet + Verifier)..." -ForegroundColor Yellow
        Run-TestList -List $IssuerTests
        Run-TestList -List $WalletTests
        Run-TestList -List $VerifierTests
    }

    $endTime = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    Write-Host "Test execution completed at $endTime" -ForegroundColor Green
    Add-Content -Path $LogFile -Value "`n=== COMPLETED at $endTime ===`n"
}
catch {
    Write-Host "Test run failed: $($_.Exception.Message)" -ForegroundColor Red
    Add-Content -Path $LogFile -Value "`n=== FAILED at $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss') ===`n$($_.Exception.Message)`n"
    exit 1
}
