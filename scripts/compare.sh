#!/bin/bash
# Compare baseline (constant ceiling) vs treatment (exponential ceiling) results.
# Parses per-SLO metrics and throughput from run output and prints a summary table.
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BUNDLE_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
RESULTS="$BUNDLE_DIR/results"

# --- Helper functions ---

parse_metric() {
  local file=$1 tier=$2 metric=$3 field=$4
  awk -v t="$tier" -v m="$metric" -v f="$field" '
    /^  [a-z]/ { cur=$1; sub(/:$/,"",cur) }
    cur==t && index($0, m":") && index($0, f"=") {
      s=$0; sub(".*"f"=", "", s); sub(/[^0-9.].*/, "", s)
      printf "%.1f\n", s/1000
    }
  ' "$file"
}

get_ttft_mean() { parse_metric "$1" "$2" "TTFT" "mean"; }
get_ttft_p99() { parse_metric "$1" "$2" "TTFT" "p99"; }
get_e2e_mean() { parse_metric "$1" "$2" "E2E" "mean"; }
get_e2e_p99() { parse_metric "$1" "$2" "E2E" "p99"; }

get_throughput() {
  local file=$1
  python3 -c "import json; d=json.load(open('$file')); print(f'{d[\"responses_per_sec\"]:.1f}')" 2>/dev/null || echo "0"
}

pct() {
  local b=$1 t=$2
  if [ -n "$b" ] && [ -n "$t" ] && awk "BEGIN{exit !($b > 0)}" 2>/dev/null; then
    awk "BEGIN { printf \"%+.0f%%\", ($t - $b) / $b * 100 }"
  else
    echo "--"
  fi
}

# --- Print results ---

printf "\n"
printf "╔══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗\n"
printf "║   Exponential Per-Band Ceiling: Baseline (constant) vs Treatment (exponential)                                                                ║\n"
printf "║   Formula: ceiling[i] = exp(-N * sat * i / (N-1))  |  Model: Qwen3-14B, trained-physics                                                      ║\n"
printf "╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝\n"

for workload in 2band_overload 2band_underload 4band_overload 4band_underload chatbot_under_4inst chatbot_mid_4inst codecompletion_under_4inst codecompletion_mid_4inst; do
  case $workload in
    2band_overload)           desc="2-band (crit 40% / batch 60%), rate=400, 5 instances [overload]" ;;
    2band_underload)          desc="2-band (crit 40% / batch 60%), rate=200, 5 instances [underload]" ;;
    4band_overload)           desc="4-band (crit 20% / std 25% / batch 30% / shed 25%), rate=400, 5 instances [overload]" ;;
    4band_underload)          desc="4-band (crit 20% / std 25% / batch 30% / shed 25%), rate=200, 5 instances [underload]" ;;
    chatbot_under_4inst)      desc="Chatbot (crit 80% / batch 20%), rate=5, 4 instances [under-capacity, 4K/1K tokens]" ;;
    chatbot_mid_4inst)        desc="Chatbot (crit 80% / batch 20%), rate=10, 4 instances [moderate, 4K/1K tokens]" ;;
    codecompletion_under_4inst) desc="CodeCompletion (crit 30% / batch 70%), rate=40, 4 instances [under-capacity, 2K/128 tokens]" ;;
    codecompletion_mid_4inst) desc="CodeCompletion (crit 30% / batch 70%), rate=95, 4 instances [moderate, 2K/128 tokens]" ;;
  esac

  # Check if results exist
  if [ ! -f "$RESULTS/baseline_${workload}_s42.txt" ]; then
    continue
  fi

  printf "\n┌──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐\n"
  printf "│ %-140s │\n" "$desc"
  printf "├──────────┬──────┬──────────────────────────────────┬──────────────────────────────────┬──────────────────────────────────┬──────────────────────────────────┤\n"
  printf "│ %-8s │ %-4s │ %-32s │ %-32s │ %-32s │ %-32s │\n" \
    "Tier" "Seed" "TTFT mean (B → T, Δ%)" "TTFT P99 (B → T, Δ%)" "E2E mean (B → T, Δ%)" "E2E P99 (B → T, Δ%)"
  printf "├──────────┼──────┼──────────────────────────────────┼──────────────────────────────────┼──────────────────────────────────┼──────────────────────────────────┤\n"

  for seed in 42 43 44; do
    baseline="$RESULTS/baseline_${workload}_s${seed}.txt"
    treatment="$RESULTS/treatment_${workload}_s${seed}.txt"

    if [ ! -f "$baseline" ] || [ ! -f "$treatment" ]; then
      continue
    fi

    for tier in critical standard batch sheddable; do
      b_ttft_mean=$(get_ttft_mean "$baseline" "$tier")
      [ -z "$b_ttft_mean" ] && continue

      t_ttft_mean=$(get_ttft_mean "$treatment" "$tier")
      b_ttft_p99=$(get_ttft_p99 "$baseline" "$tier")
      t_ttft_p99=$(get_ttft_p99 "$treatment" "$tier")
      b_e2e_mean=$(get_e2e_mean "$baseline" "$tier")
      t_e2e_mean=$(get_e2e_mean "$treatment" "$tier")
      b_e2e_p99=$(get_e2e_p99 "$baseline" "$tier")
      t_e2e_p99=$(get_e2e_p99 "$treatment" "$tier")

      d_ttft_mean=$(pct "$b_ttft_mean" "$t_ttft_mean")
      d_ttft_p99=$(pct "$b_ttft_p99" "$t_ttft_p99")
      d_e2e_mean=$(pct "$b_e2e_mean" "$t_e2e_mean")
      d_e2e_p99=$(pct "$b_e2e_p99" "$t_e2e_p99")

      col_ttft_mean=$(printf "%5s→%5sms %s" "$b_ttft_mean" "$t_ttft_mean" "$d_ttft_mean")
      col_ttft_p99=$(printf "%5s→%5sms %s" "$b_ttft_p99" "$t_ttft_p99" "$d_ttft_p99")
      col_e2e_mean=$(printf "%5s→%5sms %s" "$b_e2e_mean" "$t_e2e_mean" "$d_e2e_mean")
      col_e2e_p99=$(printf "%5s→%5sms %s" "$b_e2e_p99" "$t_e2e_p99" "$d_e2e_p99")

      printf "│ %-8s │  s%-2s │ %-32s │ %-32s │ %-32s │ %-32s │\n" \
        "$tier" "$seed" "$col_ttft_mean" "$col_ttft_p99" "$col_e2e_mean" "$col_e2e_p99"
    done
  done

  # Throughput row
  printf "├──────────┴──────┼──────────────────────────────────┴──────────────────────────────────┴──────────────────────────────────┴──────────────────────────────────┤\n"

  b_sum=0; t_sum=0; count=0
  for seed in 42 43 44; do
    bj="$RESULTS/baseline_${workload}_s${seed}.json"
    tj="$RESULTS/treatment_${workload}_s${seed}.json"
    if [ -f "$bj" ] && [ -f "$tj" ]; then
      b_t=$(get_throughput "$bj")
      t_t=$(get_throughput "$tj")
      b_sum=$(awk "BEGIN{print $b_sum + $b_t}")
      t_sum=$(awk "BEGIN{print $t_sum + $t_t}")
      count=$((count + 1))
    fi
  done

  if [ "$count" -gt 0 ]; then
    b_avg=$(awk "BEGIN{printf \"%.1f\", $b_sum/$count}")
    t_avg=$(awk "BEGIN{printf \"%.1f\", $t_sum/$count}")
    tput_d=$(pct "$b_avg" "$t_avg")
    printf "│ %-15s │ Throughput: %s → %s resp/s (%s)                                                                                              │\n" \
      "Throughput" "$b_avg" "$t_avg" "$tput_d"
  fi

  printf "└─────────────────┴───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘\n"
done

printf "\n"
printf "Legend:\n"
printf "  B = Baseline (constant ceiling 1.0, llm-d default)\n"
printf "  T = Treatment (exponential ceiling exp(-N*sat*i/(N-1)))\n"
printf "  Negative Δ%% = improvement (lower latency / less time)\n"
printf "  Positive Δ%% on batch TTFT = expected (batch held in queue to protect critical)\n"
printf "  Positive Δ%% on batch E2E  = unexpected (should investigate if large)\n"
printf "\n"
