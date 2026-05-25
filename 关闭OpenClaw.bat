@echo off

echo.
echo 正在关闭OpenClaw...
echo.

taskkill /F /IM node.exe 2>nul

echo.
echo 已关闭，可以安全拔出U盘了。
echo.
pause