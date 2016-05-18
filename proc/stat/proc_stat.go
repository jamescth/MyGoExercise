package stat

type ProcStat int

func (s ProcStat) String() string {
	return statString[s]
}

var statString []string = []string{
	"pid",
	"comm",
	"state",
	"ppid",
	"pgrp",
	"session",
	"tty_nr",
	"tpgid",
	"flags",
	"minflt",
	"cminflt",
	"majflt",
	"cmajflt",
	"utime",
	"stime",
	"cutime",
	"cstime",
	"priority",
	"nice",
	"num_threads",
	"itrealvalue",
	"starttime",
	"vsize",
	"rss",
	"rsslim",
	"startcode",
	"endcode",
	"startstack",
	"kstkesp",
	"ksteip",
	"signal",
	"blocked",
	"sigignore",
	"sigcatch",
	"wchan",
	"nswap",
	"cnswap",
	"exit_signal",
	"processor",
	"rt_priority",
	"policy",
	"delayacct_blkio_ticks",
	"guest_time",
	"cguest_time",
}

const (
	PROC_PID                   ProcStat = iota
	PROC_COMM                           //  2 The filename of the executable
	PROC_STATE                          //  3 "RSDZTW"
	PROC_PPID                           //  4 PID of the parent
	PROC_PGRP                           //  5 Group ID
	PROC_SESSION                        //  6 The session ID
	PROC_TTY_NR                         //  7 The controlling terminal of the proc
	PROC_TPGID                          //  8 The ID of the foreground proc group
	PROC_FLAGS                          //  9 The kernel flags word of proc
	PROC_MINFLT                         // 10 # of minor faults which is not page faults
	PROC_CMINFLT                        // 11 # of minor faults the proc's waited-for children have made
	PROC_MAJFLT                         // 12 # of major faults the proc has made for loading a mem page from disk
	PROC_CMAJFLT                        // 13 # of major faults the proce's wait-for children have made
	PROC_UTIME                          // 14 Amount of time the proc has been sche in user
	PROC_STIME                          // 15 Amount of time the proc has been sche in kernel
	PROC_CUTIME                         // 16 Amount of time the proc waited-for child to be sche in user
	PROC_CSTIME                         // 17 Amount of time the proc waited-for children to be sche in kernel
	PROC_PRIORITY                       // 18
	PROC_NICE                           // 19 the nice value
	PROC_NUM_THREADS                    // 20 # of threads in proc
	PROC_ITREALVALUE                    // 21 The time in jiffies b4 the next SIGALRM is sent to proc due to inter
	PROC_STARTTIME                      // 22 The time the proc started after boot
	PROC_VSIZE                          // 23 Virtual memory size in bytes
	PROC_RSS                            // 24 Resident Set Size: # of pages the proc has in real memory
	PROC_RSSLIM                         // 25 Current soft limit in bytes on the rss of the proc
	PROC_STARTCODE                      // 26 The address above which program text can run
	PROC_ENDCODE                        // 27 The address below which program text can run
	PROC_STARTSTACK                     // 28 The address of the start of the stack
	PROC_KSTKESP                        // 29 The current value of ESP
	PROC_KSTEIP                         // 30 The current EIP
	PROC_SIGNAL                         // 31 the bitmap of pending signals
	PROC_BLOCKED                        // 32 the bitmap of blocked signals
	PROC_SIGIGNORE                      // 33 The bitmap of ignored signals
	PROC_SIGCATCH                       // 34 The bitmap of caught signals
	PROC_WCHAN                          // 35 This is the "channel" in which the proc is waiting
	PROC_NSWAP                          // 36 # of pages swapped
	PROC_CNSWAP                         // 37 Cumulative nswap for child proce
	PROC_EXIT_SIGNAL                    // 38 Signal to be sent to parent when die
	PROC_PROCESSOR                      // 39 CPU # last run
	PROC_RT_PRIORITY                    // 40 Real-time scheduling priority
	PROC_POLICY                         // 41 scheduling policy
	PROC_DELAYACCT_BLKIO_TICKS          // 42 Aggregated block I/O delays
	PROC_GUEST_TIME                     // 43 Guest time of the proc
	PROC_CGUEST_TIME                    // 44 Guest time of the proc's children
)
