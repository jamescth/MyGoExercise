##### chroot example
- [* Linux / Unix: chroot Command Examples *]	(http://www.cyberciti.biz/faq/unix-linux-chroot-command-examples-usage-syntax/)

##### Network namespace
- [* Exploring networking in Linux containers *] (https://speakerdeck.com/gyre007/exploring-networking-in-linux-containers)
- [* Exploring LXC Networking *] (http://containerops.org/2013/11/19/lxc-networking/)
- [* Tenus - Golang Powered Linux Networking *] (http://containerops.org/2014/07/30/tenus-golang-powered-linux-networking/)

<br>

- [* Kernel Korner - Why and How to Use Netlink Socket *] (http://www.linuxjournal.com/article/7356)
- [* Linux Netlink as an IP Services Protocol *] (http://www.ietf.org/rfc/rfc3549.txt)

<br>

- [* How to configure a Linux bridge interface *] (http://xmodulo.com/how-to-configure-linux-bridge-interface.html)
- [* Network bridge *] (https://wiki.archlinux.org/index.php/Network_bridge)

```
	bridge using iproute2
		create:
			ip link add name bridge_james type bridge
			ip link set bridge_james up
		delete:
			ip link delete bridge_james type bridge


		veth communication:
			ip netns add james1
			ip netns add james2

			ip netns exec james1 bash
			ip link add name if_one type veth peer name if_one_peer
			ip link set dev if_one_peer netns james2
			ip addr add 192.168.0.100/24 dev if_one
			if link set if_one up

			another term:
			ip netns exec james2 bash
			ip link
			ip addr add 192.168.0.200/24 dev if_one_peer
			if link set if_one_peer up
			ping 192.168.0.100
```

<br>
- [* Introduction to IPC namespace *] (https://blog.yadutaf.fr/2013/12/28/introduction-to-linux-namespaces-part-2-ipc/)
- [* Calling setns from Go *] (http://stackoverflow.com/questions/25704661/calling-setns-from-go-returns-einval-for-mnt-namespace)

##### cgroups
- [* Linux Control Groups *] (https://sysadmincasts.com/episodes/14-introduction-to-linux-control-groups-cgroups)
- [* cgroups *] (http://lxr.free-electrons.com/source/Documentation/cgroups/cgroups.txt)

```
	cpu
		ps -o pid,psr,comm
		ps -o pid,cpuid,comm

		sudo -i
		cd /sys/fs/cgroup/cpuset
		mkdir james
		cd james
		echo 2 > cpuset.cpus
		cat cpuset.mems
		echo $$ > tasks
		ps -o pid,cpuid,comm

		rmdir james

	blkio
		1st terminal:
			sudo /usr/sbin/iotop
			sudo /usr/sbin/iotop -o

		2nd terminal:
			verify iotop
			dd if=/dev/zero of=/tmp/test-iotop bs=1M count=3000
			dd if=/dev/zero of=/tmp/test-iotop1 bs=1M count=3000
			du -h /tmp/test-iotop
				3.0G	/tmp/test-iotop
			sync
			free -m
			sudo -i
			echo 3 > /proc/sys/vm/drop_caches
			free -m
			exit

		3rd ter:
			ls -l /dev/sda*
				brw-rw---- 1 root disk 8, 0 Mar 18 08:15 /dev/sda
				brw-rw---- 1 root disk 8, 1 Mar 18 08:15 /dev/sda1
				brw-rw---- 1 root disk 8, 2 Mar 18 08:15 /dev/sda2
				brw-rw---- 1 root disk 8, 5 Mar 18 08:15 /dev/sda5

			sudo -i
			cd /sys/fs/cgroup/blkio
			mkdir test1
			cd test1
			echo "8:0 5242880" > blkio.throttle.read_bps_device 
			echo $$ > tasks 
			dd if=/tmp/test-iotop of=/dev/null
```

##### net:
- [* Using tc to priotize packets in cgroups *] (http://stackoverflow.com/questions/9904016/how-to-priotize-packets-using-tc-and-cgroups)

