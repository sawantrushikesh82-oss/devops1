param(
  [string]$OutDir = "configs"
)

$ErrorActionPreference = "Stop"

New-Item -ItemType Directory -Force -Path $OutDir | Out-Null

$privatePath = Join-Path $OutDir "jwt_private.pem"
$publicPath = Join-Path $OutDir "jwt_public.pem"
$hostKeyPath = Join-Path $OutDir "sftp_host_key"

$rsa = [System.Security.Cryptography.RSA]::Create(2048)
$privateBytes = $rsa.ExportPkcs8PrivateKey()
$privateB64 = [Convert]::ToBase64String($privateBytes)
$privatePem = "-----BEGIN PRIVATE KEY-----`n" + (($privateB64 -split '(.{1,64})' | Where-Object { $_ }) -join "`n") + "`n-----END PRIVATE KEY-----`n"
Set-Content -Path $privatePath -Value $privatePem -Encoding ascii

$publicBytes = $rsa.ExportSubjectPublicKeyInfo()
$publicB64 = [Convert]::ToBase64String($publicBytes)
$publicPem = "-----BEGIN PUBLIC KEY-----`n" + (($publicB64 -split '(.{1,64})' | Where-Object { $_ }) -join "`n") + "`n-----END PUBLIC KEY-----`n"
Set-Content -Path $publicPath -Value $publicPem -Encoding ascii

if (Get-Command ssh-keygen -ErrorAction SilentlyContinue) {
  ssh-keygen -t rsa -b 4096 -f $hostKeyPath -N "" | Out-Null
} else {
  throw "ssh-keygen was not found on PATH. Install OpenSSH Client and re-run."
}

Write-Host "Generated:"
Write-Host " - $privatePath"
Write-Host " - $publicPath"
Write-Host " - $hostKeyPath"
