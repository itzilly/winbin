@echo off
setlocal ENABLEDELAYEDEXPANSION

if not exist bin (
    mkdir bin
)

for /D %%f in (cmd\*) do (
    set "tool=%%~nxf"
    echo Building %%f...
    go build -o "bin\!tool!.exe" ".\cmd\!tool!"
)

echo Done!
