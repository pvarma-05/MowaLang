@echo off
set INSTALL_DIR=%ProgramFiles%\mowalang
set BIN_URL=https://mowalang.vercel.app/bin/windows/mowa.exe

REM Check if the "mowalang" folder exists
if not exist "%INSTALL_DIR%" (
    echo Creating directory "%INSTALL_DIR%"...
    mkdir "%INSTALL_DIR%"
) else (
    echo "%INSTALL_DIR%" already exists. Skipping directory creation.
)

REM Navigate to the directory
cd /d "%INSTALL_DIR%"

REM Check if mowa.exe exists, delete if it does
if exist "mowa.exe" (
    echo Deleting existing mowa.exe...
    del /f mowa.exe
)

REM Download the latest mowa.exe using curl (or wget)
echo Downloading mowa.exe...
curl -L %BIN_URL% -o mowa.exe

REM Check if the path is in the system environment variable, if not, add it
setlocal enabledelayedexpansion
set PATH_EXISTS=0
for %%A in ("%PATH%") do (
    if "%%~A"=="%INSTALL_DIR%" set PATH_EXISTS=1
)

if %PATH_EXISTS%==0 (
    echo Adding "%INSTALL_DIR%" to system PATH...
    setx PATH "%PATH%;%INSTALL_DIR%"
)

REM Completion message
echo Super mowa! nuvvu mowalang successfull ga install chesaav..
pause
