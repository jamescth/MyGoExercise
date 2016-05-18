namespaces(7)
	overview of Linux namespaces

clone(2)
	create a child process
		CLONE_NEWIPC - new IPC namespace (Linux 2.6.19)
			svipc(7), mq_overview(7)
		CLONE_NEWNET - new network namespace (Linux 2.6.24)
			/proc/net and /sys/class/net, namespaces(7)
		CLONE_NEWNS  - new mount namespace (Linux 2.4.19)
			namespaces(7)
		CLONE_NEWPID - new PID namespace (Linux 2.6.24)
			namespaces(7), pid_namespaces(7)
		CLONE_NEWUSER - new user namespace (Linux 3.8)
			namespaces(7), user_namespaces(7)
		CLONE_NEWUTS - new UTS namespace (Linux 2.6.19)
			uname(2), setdomainname(2), sethostname(2)

unshare(2)
	disassociate parts of the process execution context

setns(2)
	reassociate thread with a namespace

pivot_root(2)
	change the root filesystem, used in boot time.
	*** may not suit for namespace ***

chroot(2)
	change root directory

switch_root(8)

************************************
cgroups
http://lxr.free-electrons.com/source/Documentation/cgroups/cgroups.txt?v=3.15

systemd(1)
systemd.cgroup(5)
systemd.mount(5)

cgconfig.conf(5) - libcgroup configuration file
