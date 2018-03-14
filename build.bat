@echo off

IF NOT EXIST build mkdir build

pushd src
go build -o ..\build\render.exe

..\build\render.exe test.png
popd