#!/bin/bash

# Script de stress test con diferentes n
for n in 100 500 1000 2000 5000; do
  for route in "on" "onlogn" "on2"; do
    for i in $(seq 1 20); do
      curl -s "https://powerbi-rest-production.up.railway.app/$route?n=$n" > /dev/null
    done
  done
done

# O(2^n) con valores pequeños para no explotar el CPU
for n in 10 15 20 25 30 35; do
  for i in $(seq 1 15); do
    curl -s "https://powerbi-rest-production.up.railway.app/o2n?n=$n" > /dev/null
  done
done
