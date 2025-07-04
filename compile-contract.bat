@echo off
echo Compiling Neo smart contract...
neo-go contract compile -i main.go -c main.yml -m main.manifest.json -o main.nef
echo Compilation complete!
pause
