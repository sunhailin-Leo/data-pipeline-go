#!/bin/bash
set -e

echo "Running benchmarks..."

BENCH_PACKAGES=(
  "./pkg/utils/..."
  "./pkg/transform/convert/..."
  "./pkg/middlewares/..."
)

OUTPUT_FILE="benchmark_results.txt"

> "$OUTPUT_FILE"

for pkg in "${BENCH_PACKAGES[@]}"; do
  echo "Benchmarking $pkg ..."
  go test -bench=. -benchmem -count=5 -run=^$ "$pkg" | tee -a "$OUTPUT_FILE"
done

echo ""
echo "Results saved to $OUTPUT_FILE"

if command -v benchstat &> /dev/null; then
  echo ""
  echo "=== Benchmark Statistics ==="
  benchstat "$OUTPUT_FILE"
fi
