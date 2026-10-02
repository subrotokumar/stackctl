# 1. Define installation directory
$targetDir = "$env:LOCALAPPDATA\Programs\stackctl"
New-Item -ItemType Directory -Force -Path $targetDir | Out-Null

# 2. Detect System Architecture (x86_64 vs arm64)
$arch = if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "x86_64" }

# 3. Fetch the latest release tag from GitHub API
$tag = (Invoke-RestMethod "https://api.github.com/repos/subrotokumar/stackctl/releases/latest").tag_name

# 4. Download release archive based on detected architecture
curl.exe -sL "https://github.com/subrotokumar/stackctl/releases/download/$tag/stackctl_Windows_$arch.zip" -o "$targetDir\stackctl.zip"

# 5. Extract executable and cleanup zip
Expand-Archive -Path "$targetDir\stackctl.zip" -DestinationPath $targetDir -Force
Remove-Item "$targetDir\stackctl.zip"

# 6. Permanently add target directory to User PATH
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$targetDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$targetDir", "User")
}

# 7. Refresh PATH for current active session
$env:Path += ";$targetDir"