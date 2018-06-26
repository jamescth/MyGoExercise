package config

import "time"

type LogFile struct {
	File          string
	LinkTo        string
	Type          string
	StartTime     time.Time
	EndTime       time.Time
	FirstLineTime time.Time
	LastLineTime  time.Time
}

/*
upgrade milestone
{SRC}/UpgradeMgr/IDL/UpgradeTaskInfo.proto

enum UpgradeCheckpoint {
16    TASK_STARTING                  = 0;
17    TASK_STARTED                   = 1;
18    CHECK_PASSED                   = 2;
19    IMAGE_DOWNLOADED               = 3;
20    IMAGE_DEPLOYED                 = 4;
21    NONDISRUPT_STARTED             = 100;
22    AGENTS_DOWNLOAD_SET            = 101;
23    AGENTS_DOWNLOAD_CONVERGED      = 102;
24    AGENTS_PREPARE_SET             = 103;
25    AGENTS_PREPARE_CONVERGED       = 104;
26    AGENTS_READY                   = 105;
27    HOSTS_SWITCHOVER_CALLED        = 106;
28    HOSTS_SWITCHOVER_WAITED        = 107;
29    OTHER_NODES_SWITCHOVER_CALLED  = 108;
30    DISRUPTIVE_STARTED             = 150;
31    VERSION_SWITCHING              = 200; // Do not change this value
32    VERSION_SWITCHED               = 201;
33    VCPLUGIN_REREGISTERED          = 210;
34    TASK_FINISHED                  = 250;
35    TASK_FAILURE                   = 251;
36    TASK_THE_END                   = 300;
37}
*/
