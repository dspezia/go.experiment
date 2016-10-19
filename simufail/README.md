# simufail

This project implements a Monte-Carlo simulation to calculate the probability of failures of a Couchbase cluster.

The idea is to parameter a number of random failures in a one year timeframe, and then simulate the cluster many times, in order to deduce the probability to experience various issues in the year.

To use it, edit the config.go file to alter the parameters.
Build it, and run the simufail binary.

    Usage of ./simufail:
      -b int
        	Size of batch (default 100000)
      -n int
        	Number of iterations (default 32)
      -p int
        	Parallelism level (default 8)
      -s int
        	Seed for random numbers

## Cluster topology and failure events

We consider a cluster with 3 replicas (i.e. 1 master and 2 slaves per vbucket) and 3 server groups, with rack awareness.

### Unexpected events:

- Each node suffers from N_REBOOTS unexpected reboot a year, resulting in MTBR_REBOOT secs outage per node.
- Each node has a PROB_HW_FAILURE% chance a year to get a hardware failure resulting in MTBR_HW_FAILURE secs outage.
- Each availability zone has a PROB_NET_FAILURE% chance a year to suffer from a network issue making it unavailable for MTBR_NET_FAILURE secs.

### Scheduled events:

- Each availability zone is brought down N_ZONE_SHUTDOWNS times a year resulting in MTBR_ZONE_SHUTDOWN secs outage per zone.
- The cluster is upgraded N_ROLLING_UPGRADES times a year using rolling upgrade, resulting in MTBR_ROLLING_UPGRADE secs outage per node in sequence, separated by IDLE_ROLLING_UPGRADE secs idle periods.

Availability zone shutdowns and couchbase rolling upgrades are scheduled events, so:
 - they are mutually exclusive
 - they are not scheduled when there is already a node down for any reason
