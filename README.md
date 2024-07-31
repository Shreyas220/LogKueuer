# LogKueuer

This is an experiment to develop a Kubernetes-native computational engine.

The goal is 
- **Data Partitioning:** Segment log data into manageable chunks 
- **Distributed Job Execution:** Deploy Kubernetes jobs to process each data chunk concurrently.
- **Computation:** Do computational tasks on the distributed log chunks.
- **Result Aggregation:** Consolidate the results from multiple Kubernetes jobs.
