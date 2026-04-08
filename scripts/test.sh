#!/usr/bin/env bash

set -euo pipefail

API_BASE="${API_BASE:-https://powerbi-rest-production.up.railway.app}"
REQUESTS_PER_CASE="${REQUESTS_PER_CASE:-5}"
BENCHMARK_RUNS="${BENCHMARK_RUNS:-3}"

N_VALUES=(100000 500000 1000000)
PARALLEL_WORKERS=(2 4 8)
THREAD_WORKERS=(4 8 12)

run_case() {
  local route="$1"
  local query="$2"
  local url="${API_BASE}/${route}"

  if [[ -n "${query}" ]]; then
    url="${url}?${query}"
  fi

  for _ in $(seq 1 "${REQUESTS_PER_CASE}"); do
    curl -fsS --max-time 30 "${url}" > /dev/null
  done

  printf "OK  %s  x%s\n" "${url}" "${REQUESTS_PER_CASE}"
}

echo "API base: ${API_BASE}"
echo "Iniciando pruebas de rutas paralelas..."

for n in "${N_VALUES[@]}"; do
  run_case "sequential" "n=${n}&runs=${BENCHMARK_RUNS}"
done

for n in "${N_VALUES[@]}"; do
  for workers in "${PARALLEL_WORKERS[@]}"; do
    run_case "parallel" "n=${n}&runs=${BENCHMARK_RUNS}&workers=${workers}"
  done
done

for n in "${N_VALUES[@]}"; do
  for workers in "${THREAD_WORKERS[@]}"; do
    run_case "parallel-with-threads" "n=${n}&runs=${BENCHMARK_RUNS}&workers=${workers}"
  done
done

for n in "${N_VALUES[@]}"; do
  for p_workers in "${PARALLEL_WORKERS[@]}"; do
    for t_workers in "${THREAD_WORKERS[@]}"; do
      run_case "parallel-metrics" "n=${n}&runs=${BENCHMARK_RUNS}&parallel_workers=${p_workers}&thread_workers=${t_workers}"
    done
  done
done

run_case "summary" ""
echo "Pruebas completadas."
