# Configuration

## llm-d Real Deployment

To reproduce the experiment on a real llm-d cluster, deploy with this setup.

### vLLM Pod Configuration

| Parameter | Value | Notes |
|---|---|---|
| Model | `Qwen/Qwen3-14B` | |
| GPU | H100-SXM-80GB | |
| `tensor_parallel_size` | 1 | |
| `max_num_seqs` | 256 | Max concurrent requests per pod |
| `max_num_batched_tokens` | 2048 | Chunked prefill budget |
| `block_size` | 16 | KV cache block size in tokens |
| `gpu_memory_utilization` | 0.9 | |
| `max_model_len` | 40960 | |
| `enable_chunked_prefill` | True | |
| Number of pods | 4 | |

### llm-d EPP Configuration (Baseline)

Flow control enabled with default constant ceiling. Use `random-picker` for endpoint selection:

```yaml
apiVersion: inference.networking.x-k8s.io/v1alpha1
kind: EndpointPickerConfig
featureGates:
  - flowControl
plugins:
  - type: random-picker
schedulingProfiles:
  - name: default
    plugins:
      - pluginRef: random-picker
flowControl: {}
```

### llm-d EPP Configuration (Treatment)

Add the exponential ceiling plugin (see README for full Go code):

```yaml
apiVersion: inference.networking.x-k8s.io/v1alpha1
kind: EndpointPickerConfig
featureGates:
  - flowControl
plugins:
  - type: exponential-ceiling-policy
    name: exponential-ceiling
  - type: random-picker
schedulingProfiles:
  - name: default
    plugins:
      - pluginRef: random-picker
flowControl:
  usageLimitPolicyPluginRef: "exponential-ceiling"
```

### Priority Bands (InferenceObjective CRs)

Two tiers used across all workloads:

| Tier | Priority value | Role |
|---|---|---|
| critical | 100 | Protected — ceiling stays at 1.0 under load |
| sheddable | -50 | Low-priority — gated earliest, held in queue |

```yaml
apiVersion: inference.networking.x-k8s.io/v1alpha2
kind: InferenceObjective
metadata:
  name: critical
spec:
  priority: 100
---
apiVersion: inference.networking.x-k8s.io/v1alpha2
kind: InferenceObjective
metadata:
  name: sheddable
spec:
  priority: -50
```

---

## BLIS Simulation (optional — for reproducing the scripts)

The scripts (`scripts/run.sh`) use these BLIS flags:

| Flag | Value | Notes |
|---|---|---|
| `--model` | qwen/qwen3-14b | |
| `--latency-model` | trained-physics | Physics-informed with learned corrections |
| `--max-model-len` | 40960 | |
| `--flow-control` | (enabled) | |
| `--saturation-detector` | utilization | |
| `--queue-depth-threshold` | 5 | |
| `--kv-cache-util-threshold` | 0.8 | |
| `--dispatch-order` | priority | |
| `--usage-limit-threshold` | 1.0 | |
| `--num-instances` | 4 | |

BLIS defaults used (not overridden, auto-derived from model/hardware):

| Parameter | Value | Corresponds to vLLM |
|---|---|---|
| `--max-num-running-reqs` | 256 | `max_num_seqs=256` |
| `--max-num-scheduled-tokens` | 2048 | `max_num_batched_tokens=2048` |
| `--block-size-in-tokens` | 16 | `block_size=16` |
| `--gpu-memory-utilization` | 0.9 | `gpu_memory_utilization=0.9` |
| `--total-kv-blocks` | auto-calculated | Derived from model size, GPU memory, gpu_memory_utilization |

Treatment is applied via `scripts/treatment.patch` (one-line formula change in `DequeueGated()`).

## Commits

| Component | Commit | Notes |
|---|---|---|
| BLIS (scripts) | `0195abf` | Latest main (includes PR #1382 fix) |
| Nous (campaign) | `d6245af` | main |
| llm-d (code proofs) | `a10f9ac8` | Source references in README |
