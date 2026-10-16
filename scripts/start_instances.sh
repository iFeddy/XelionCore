#!/usr/bin/env bash
# Script to launch all 3 Fiesta Online server instances in the background

echo "Starting XelionCore server instances..."

# Kill any existing instances
pkill -f "./bin/login" 2>/dev/null
pkill -f "./bin/world" 2>/dev/null
pkill -f "./bin/zone" 2>/dev/null
pkill -f "./bin/server" 2>/dev/null

./bin/login > logs/login.log 2>&1 &
echo "[1/3] Login Server started on 0.0.0.0:9010 (PID $!)"

./bin/world > logs/world.log 2>&1 &
echo "[2/3] World/Character Server started on 0.0.0.0:9110 (PID $!)"

./bin/zone > logs/zone.log 2>&1 &
echo "[3/3] Zone Server started on 0.0.0.0:9210 (PID $!)"

echo "All 3 instances running! Logs available in logs/."
