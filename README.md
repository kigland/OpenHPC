# HPC Scheduler

A lightweight workload scheduler. Similar to Slurm, but simpler and easier for small-size distributed HPC cluster.

```mermaid
flowchart TB
  M[/Coordinator\]
  P[/Proxier\]
  U([Users])
  subgraph HPC Cluster
    HPC1[HPC 1]
    HPC2[HPC 2]
    HPC3[HPC 3]
  end
  M --> HPC1
  M --> HPC2
  M --> HPC3
  P --> HPC1
  P --> HPC2
  P --> HPC3
  M --> P
  U -->|Submit| M
  U -->|Access| P
```