@echo off

set "FOUND_DRIVE="

for %%D in (D E F G H I J K L M N O P Q R S T U V W X Y Z) do (
    if exist "%%D:\openclaw\" (
        set "FOUND_DRIVE=%%D:"
        goto :found_drive
    )
)

:found_drive
if "%FOUND_DRIVE%"=="" (
    echo [错误] 未找到OpenClaw！
    echo 请确认已解压到U盘根目录。
    pause
    exit /b 1
)

set "USB_ROOT=%FOUND_DRIVE%"
set "PATH=%USB_ROOT%\node;%PATH%"

cd /d "%USB_ROOT%"
start "OpenClaw" "%USB_ROOT%\node\node.exe" "%USB_ROOT%\setup.js"

pause