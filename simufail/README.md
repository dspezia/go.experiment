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

## Results

Here are some results for 12-nodes, 9-nodes and 6-nodes clusters distributed over 3 availability zones.
Generated with:

    N_ZONES              = 3
    N_ROLLING_UPGRADES   = 2
    MTBR_ROLLING_UPGRADE = 2 * 60
    IDLE_ROLLING_UPGRADE = 30
    N_ZONE_SHUTDOWNS     = 1
    MTBR_ZONE_SHUTDOWN   = 3 * 3600
    N_REBOOTS            = 3
    MTBR_REBOOT          = 5 * 60
    PROB_NET_FAILURE     = 25
    MTBR_NET_FAILURE     = 3600
    PROB_HW_FAILURE      = 6
    MTBR_HW_FAILURE      = 3 * 24 * 3600
    THROUGHPUT           = 50

### For a 12 nodes cluster

    Starting 256 iterations over 8 threads with 100000 batch size
    Iteration 255
    Done

    Number of simulations:            25600000
    Average node failures per year:   64.20

    Failures on 3 zones

    Average % of keys impacted:      1.9814 %
    Average duration of the event:   1329 secs
    Average number of records lost:  0.9907

    Histogram of probability
    At least  1 occurences:    0.0164 %
    At least  2 occurences:    0.0005 %

    Failures on 2 zones

    Average % of keys impacted:      9.5534 %
    Average duration of the event:   2264 secs
    Average number of transactions:  10815.9078

    Histogram of probability
    At least  1 occurences:   15.9393 %
    At least  2 occurences:    2.5502 %
    At least  3 occurences:    0.4583 %
    At least  4 occurences:    0.0664 %
    At least  5 occurences:    0.0144 %
    At least  6 occurences:    0.0070 %
    At least  7 occurences:    0.0049 %
    At least  8 occurences:    0.0033 %
    At least  9 occurences:    0.0008 %
    At least 10 occurences:    0.0001 %
    More occurences:           0.0000 %

    Failures involving at least 2 nodes

    Average % of keys impacted:      94.9560 %
    Average duration of the event:   8861 secs
    Average number of transactions:  420680.0638

    Histogram of probability
    At least  1 occurences:   65.9787 %
    At least  2 occurences:   26.4358 %
    At least  3 occurences:    7.5125 %
    At least  4 occurences:    1.8742 %
    At least  5 occurences:    0.4282 %
    At least  6 occurences:    0.0954 %
    At least  7 occurences:    0.0269 %
    At least  8 occurences:    0.0123 %
    At least  9 occurences:    0.0087 %
    At least 10 occurences:    0.0067 %
    More occurences:           0.0052 %

### For a 9 nodes cluster

    Starting 256 iterations over 8 threads with 100000 batch size
    Iteration 255
    Done

    Number of simulations:            25600000
    Average node failures per year:   49.13

    Failures on 3 zones

    Average % of keys impacted:      4.6417 %
    Average duration of the event:   1621 secs
    Average number of records lost:  2.3209

    Histogram of probability
    At least  1 occurences:    0.0070 %
    At least  2 occurences:    0.0001 %
    At least  3 occurences:    0.0000 %

    Failures on 2 zones

    Average % of keys impacted:      16.0693 %
    Average duration of the event:   2190 secs
    Average number of transactions:  17592.2673

    Histogram of probability
    At least  1 occurences:   10.0041 %
    At least  2 occurences:    1.1053 %
    At least  3 occurences:    0.1621 %
    At least  4 occurences:    0.0178 %
    At least  5 occurences:    0.0046 %
    At least  6 occurences:    0.0025 %
    At least  7 occurences:    0.0003 %
    At least  8 occurences:    0.0000 %
    At least  9 occurences:    0.0000 %

    Failures involving at least 2 nodes

    Average % of keys impacted:      97.3129 %
    Average duration of the event:   9053 secs
    Average number of transactions:  440475.7563

    Histogram of probability
    At least  1 occurences:   62.7979 %
    At least  2 occurences:   22.1077 %
    At least  3 occurences:    4.9543 %
    At least  4 occurences:    0.9447 %
    At least  5 occurences:    0.1721 %
    At least  6 occurences:    0.0307 %
    At least  7 occurences:    0.0092 %
    At least  8 occurences:    0.0047 %
    At least  9 occurences:    0.0031 %
    At least 10 occurences:    0.0023 %
    More occurences:           0.0018 %

### For a 6 nodes cluster

    Starting 256 iterations over 8 threads with 100000 batch size
    Iteration 255
    Done

    Number of simulations:            25600000
    Average node failures per year:   34.04

    Failures on 3 zones

    Average % of keys impacted:      14.2829 %
    Average duration of the event:   1217 secs
    Average number of records lost:  7.1414

    Histogram of probability
    At least  1 occurences:    0.0024 %
    At least  2 occurences:    0.0000 %

    Failures on 2 zones

    Average % of keys impacted:      32.6166 %
    Average duration of the event:   2056 secs
    Average number of transactions:  33536.2500

    Histogram of probability
    At least  1 occurences:    5.1818 %
    At least  2 occurences:    0.3400 %
    At least  3 occurences:    0.0361 %
    At least  4 occurences:    0.0031 %
    At least  5 occurences:    0.0002 %
    At least  6 occurences:    0.0000 %

    Failures involving at least 2 nodes

    Average % of keys impacted:      99.0190 %
    Average duration of the event:   9203 secs
    Average number of transactions:  455620.6947

    Histogram of probability
    At least  1 occurences:   60.2276 %
    At least  2 occurences:   18.7979 %
    At least  3 occurences:    3.1990 %
    At least  4 occurences:    0.4184 %
    At least  5 occurences:    0.0588 %
    At least  6 occurences:    0.0077 %
    At least  7 occurences:    0.0023 %
    At least  8 occurences:    0.0013 %
    At least  9 occurences:    0.0009 %
    At least 10 occurences:    0.0003 %
    More occurences:           0.0001 %
