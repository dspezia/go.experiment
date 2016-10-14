package main

// Unexpected events:
// Each node suffers from N_REBOOTS unexpected reboot a year, resulting in MTBR_REBOOT secs outage per node.
// Each node has a PROB_HW_FAILURE% chance a year to get a hardware failure resulting in MTBR_HW_FAILURE secs outage.
// Each availability zone has a PROB_NET_FAILURE% chance a year to suffer from a network issue making it unavailable for MTBR_NET_FAILURE secs.

// Scheduled events:
// Each availability zone is brought down N_ZONE_SHUTDOWNS times a year resulting in MTBR_ZONE_SHUTDOWN secs outage per zone.
// The cluster is upgraded N_ROLLING_UPGRADES times a year using rolling upgrade, resulting in MTBR_ROLLING_UPGRADE secs outage
// per node in sequence, separated by IDLE_ROLLING_UPGRADE secs idle periods.
// Availability zone shutdowns and couchbase rolling upgrades are scheduled events, so:
//   - they are mutually exclusive
//   - they are not scheduled when there is already a node down for any reason
const (
	N_NODES              = 12
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
	ZONE_SIZE            = N_NODES / N_ZONES
)
