**Run Summary: `r1`**
Generated: 2026-05-20T17:16:10.047290+00:00 | Scenario: exponential-ceiling

**Algorithm**
- Source: `algorithms/exponential_ceiling.go`
- Description: Exponential per-band dispatch ceiling policy — ceiling[i] = exp(-N * sat * i / (N-1))

**Translation**
- Plugin type: `exponential-ceiling-policy`
- Files created: `epp`, `pkg/epp/framework/plugins/flowcontrol/usagelimits/exponentialceiling/policy.go`, `pkg/epp/framework/plugins/flowcontrol/usagelimits/exponentialceiling/policy_test.go`
- Files modified: `cmd/epp/runner/runner.go`
- Review: 1/1 after 1 rounds

**Packages**

- `/Users/kalantar/projects/go.workspace/src/github.com/kalantar-msb/exponential-ceiling/workspace/runs/r1/cluster/pipelinerun-balanced-mid-baseline1.yaml`
- `/Users/kalantar/projects/go.workspace/src/github.com/kalantar-msb/exponential-ceiling/workspace/runs/r1/cluster/pipelinerun-balanced-mid-treatment.yaml`
- `/Users/kalantar/projects/go.workspace/src/github.com/kalantar-msb/exponential-ceiling/workspace/runs/r1/cluster/pipelinerun-balanced-under-baseline1.yaml`
- `/Users/kalantar/projects/go.workspace/src/github.com/kalantar-msb/exponential-ceiling/workspace/runs/r1/cluster/pipelinerun-balanced-under-treatment.yaml`
- `/Users/kalantar/projects/go.workspace/src/github.com/kalantar-msb/exponential-ceiling/workspace/runs/r1/cluster/pipelinerun-blindspot-mid-baseline1.yaml`
- `/Users/kalantar/projects/go.workspace/src/github.com/kalantar-msb/exponential-ceiling/workspace/runs/r1/cluster/pipelinerun-blindspot-mid-treatment.yaml`
- `/Users/kalantar/projects/go.workspace/src/github.com/kalantar-msb/exponential-ceiling/workspace/runs/r1/cluster/pipelinerun-blindspot-under-baseline1.yaml`
- `/Users/kalantar/projects/go.workspace/src/github.com/kalantar-msb/exponential-ceiling/workspace/runs/r1/cluster/pipelinerun-blindspot-under-treatment.yaml`
- `/Users/kalantar/projects/go.workspace/src/github.com/kalantar-msb/exponential-ceiling/workspace/runs/r1/cluster/pipelinerun-chatbot-mid-baseline1.yaml`
- `/Users/kalantar/projects/go.workspace/src/github.com/kalantar-msb/exponential-ceiling/workspace/runs/r1/cluster/pipelinerun-chatbot-mid-treatment.yaml`
- `/Users/kalantar/projects/go.workspace/src/github.com/kalantar-msb/exponential-ceiling/workspace/runs/r1/cluster/pipelinerun-chatbot-under-baseline1.yaml`
- `/Users/kalantar/projects/go.workspace/src/github.com/kalantar-msb/exponential-ceiling/workspace/runs/r1/cluster/pipelinerun-chatbot-under-treatment.yaml`
- `/Users/kalantar/projects/go.workspace/src/github.com/kalantar-msb/exponential-ceiling/workspace/runs/r1/cluster/pipelinerun-codecompletion-mid-baseline1.yaml`
- `/Users/kalantar/projects/go.workspace/src/github.com/kalantar-msb/exponential-ceiling/workspace/runs/r1/cluster/pipelinerun-codecompletion-mid-treatment.yaml`
- `/Users/kalantar/projects/go.workspace/src/github.com/kalantar-msb/exponential-ceiling/workspace/runs/r1/cluster/pipelinerun-codecompletion-under-baseline1.yaml`
- `/Users/kalantar/projects/go.workspace/src/github.com/kalantar-msb/exponential-ceiling/workspace/runs/r1/cluster/pipelinerun-codecompletion-under-treatment.yaml`

**Workloads**

- balanced_mid
- balanced_under
- blindspot_mid
- blindspot_under
- chatbot_mid
- chatbot_under
- codecompletion_mid
- codecompletion_under

**Checklist**
- [x] Translation complete
- [x] Assembly complete
- [x] validate-assembly passed

**Verdict: READY TO DEPLOY**
