# Kubewarden Runtime Enforcer

Kubewarden Runtime Enforcer is a Kubernetes security tool that uses eBPF
(Extended Berkeley Packet Filter) to observe process executions in your
workloads and to enforce allow-list based security policies at the kernel
level.

It operates in three phases:

- **Learn** — observe process executions and generate a `WorkloadPolicyProposal`
  per workload.
- **Monitor** — report violations of an approved `WorkloadPolicy` without
  blocking them.
- **Protect** — block executions that violate the policy allow-list.

The project is in Beta. The design documents are in the [RFCs](docs/rfc).

## Documentation

The full documentation is available at
[docs.kubewarden.io/runtime-enforcer](https://docs.kubewarden.io/runtime-enforcer/latest/en/introduction.html).

- [Kubewarden Runtime Enforcer Quick Start](https://docs.kubewarden.io/runtime-enforcer/latest/en/installation/quickstart.html)
  — deploy Runtime Enforcer and walk through the learn/monitor/protect workflow.
- [Compatibility](https://docs.kubewarden.io/runtime-enforcer/latest/en/compatibility.html)
  — kernel, architecture and platform requirements.
- [Kubewarden Runtime Enforcer phases: learn, monitor, protect](https://docs.kubewarden.io/runtime-enforcer/latest/en/phases.html)
  — understand the learn, monitor and protect phases in detail.

## Build with the community

Runtime Enforcer is part of [Kubewarden](https://www.kubewarden.io/), a CNCF
Sandbox project. Rancher by SUSE developed it first.

Get help, share an idea, or improve these docs.

- [Join us on Slack](https://kubernetes.slack.com/?redir=%2Fmessages%2Fkubewarden)
- [Community meetings](https://www.kubewarden.io/#get-in-touch)
- [Contribute to the docs](https://github.com/kubewarden/docs)
- [Security policy](https://github.com/kubewarden/community/security/policy)
