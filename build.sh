#!/bin/bash
clear 
echo "BUILDING EMERALD"
sleep 1
go build -o emerald cmd/main.go
sleep 1
echo "BUILDED EMERALD!!!! Run: ./emerald"