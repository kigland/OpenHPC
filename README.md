# Kigland OpenHPC

![](docs/assets/portal.png)

A lightweight workload scheduler. Similar to Slurm, but simpler and easier for small-size distributed HPC cluster.

## Overview

What we have done:

- Multi GPU Assignment
- Safe Environment environment for different users
- Multi port access (HTTP, SSH, ...) (Simple Proxier)
- Web Portal

TODO:

- Job Queue
- Job Scheduler (SJF/FCFS)
- User Portal
- Multiple Physical Node Coordination
- ...

## Architecture

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