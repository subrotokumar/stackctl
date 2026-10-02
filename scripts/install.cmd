@echo off
setlocal enabledelayedexpansion

:: 1. Define installation directory (%LOCALAPPDATA%\Programs\stackctl)
set "TARGET_DIR=%LOCALAPPDATA%\Programs\stackctl"
if not exist "%TARGET_DIR%" mkdir "%TARGET_DIR%"

:: 2. Detect System Architecture (x86_64 vs arm64)
if /i "%PROCESSOR_ARCHITECTURE%"=="ARM64" (
    set "ARCH=arm64"
) else (
    set "ARCH=x86_64"
)

:: 3. Fetch the latest release tag from GitHub API using PowerShell
for /f "usebackq tokens=*" %%A in (`powershell -Command "(Invoke-RestMethod 'https://api.github.com/repos/subrotokumar/stackctl/releases/latest').tag_name"`) do set "TAG=%%A"

:: 4. Download release archive based on detected architecture
curl -sL "https://github.com/subrotokumar/stackctl/releases/download/%TAG%/stackctl_Windows_%ARCH%.zip" -o "%TARGET_DIR%\stackctl.zip"

:: 5. Extract archive using tar and delete zip file
tar -xf "%TARGET_DIR%\stackctl.zip" -C "%TARGET_DIR%"
del "%TARGET_DIR%\stackctl.zip"

:: 6. Add to User PATH if not already present
echo ;%PATH%; | find /i ";%TARGET_DIR%;" >nul
if errorlevel 1 (
    for /f "tokens=2*" %%A in ('reg query "HKCU\Environment" /v Path 2^>nul') do set "USER_PATH=%%B"
    if defined USER_PATH (
        setx Path "%USER_PATH%;%TARGET_DIR%" >nul
    ) else (
        setx Path "%TARGET_DIR%" >nul
    )
)

:: 7. Refresh PATH in current CMD session
set "PATH=%PATH%;%TARGET_DIR%"

echo stackctl (%ARCH%) installed successfully to %TARGET_DIR%