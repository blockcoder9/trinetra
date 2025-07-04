@echo off
echo Running Neo RPC Stability Tests...
echo.

go run rpc-stability-test.go

echo.
echo Test completed! Check the results above for RPC endpoint recommendations.
pause
