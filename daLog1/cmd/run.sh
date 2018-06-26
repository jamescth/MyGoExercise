#!/bin/bash
#./daLog -path /cores/bug_30350
#./daLog -path /cores/bug_29487

#./daLog -is_http 1 
# cat result.out | tr -d '\000' | grep -e ABORT -e DA_ASSERT -e dafile -e checkpoint -e exception -e CalledProcessError -e "Caught exception"
# ABORT/DA_ASSERT/DA_PANIC: core dump
# dumping: esx vmkernel show core dump
# checkpoint / Handling : Upgrade
# Caught execption : python exception 
# lingering : linger processes procmgr
# Handling RpcUpgradeAgent: Upgrade agent state
# Handling RpcUpgradeMaster: Upgrade master
# checkpoint= : Upgrade Server state
# DaExceptionRpcTimeout: 
# Killing UpgradeMgr: test.log to kill UpgradeMgr in nondisruptive
# Failover: hamgr
# TODO: events show:
# 
# received req checkUpgradeState / check_upgrade_state : controller01/da/data/var/log/hamgr.log
#	Controller(1, role=Active, state=Ok, mirroring=False, resyncing=True).check_upgrade_state: running=4.0.101.0-30762_93c05b8_u_g, prepare=9.0.8.0-999999_abcdef, mirror=False, peer_running=4.0.101.0-30762_93c05b8_u_g
#	Controller(1, role=Passive, state=Resyncing, mirroring=False, resyncing=False).upgrade_and_reboot_if_needed: Upgrading from 4.0.101.0-30762_93c05b8_u_g to 9.0.8.0-999999_abcdef
#	[DSBPoller] images on [(u'4.0.101.0-30762_93c05b8_u_g', 975585280, u'2018-03-27 15:31:34', None, 1522199287.305035, 0, '/dev/sw_image01'), (u'4.0.101.0-30762_93c05b8_u_g', 975585280, u'2018-03-27 15:31:34', None, 1522199292.476694, 0, '/dev/sw_image02'), (u'4.0.101.0-30762_93c05b8_u_g', 975585280, u'2018-03-27 15:31:34', None, 1522199300.737376, 0, '/dev/sw_image03'), (u'9.0.8.0-999999_abcdef', 975575040, u'2018-03-27 15:11:08', u'https://upgrade-center-test.datrium.com:443/release_note?requested_version=9.0.8.0-999999_abcdef', 1522200517.855051, 1, '/dev/sw_image04'), (u'9.0.8.0-999999_abcdef', 975575040, u'2018-03-27 15:11:08', u'https://upgrade-center-test.datrium.com:443/release_note?requested_version=9.0.8.0-999999_abcdef', 1522200525.003078, 1, '/dev/sw_image05'), (u'9.0.8.0-999999_abcdef', 975575040, u'2018-03-27 15:11:08', u'https://upgrade-center-test.datrium.com:443/release_note?requested_version=9.0.8.0-999999_abcdef', 1522200531.412952, 1, '/dev/sw_image06')]
#                                  {controller-VMware-421f36113c0cd431-6b816fbc48a3a5c6} {img_store.py:492:_list_files}
#
# controller procmgr.log
#  starting process UpgradeMgr / started process UpgradeMgr / killing process group UpgradeMgr / process UpgradeMgr (773) terminated, coredump: False, signal: 9
#
# controller upgrade_mgr
#  2018-03-28T01:08:34.025 DA_LOG_INFO [UpgradeCompareVersions:ace83c] Agent(node1.controller1, ctrl-5aba2644-0211-45dd-a22d-8391cad4748f): ZK currentStates updated: RUN [4.0.101.0-30762_93c05b8_u_g] DL [none] PREP [none]
# Upgrade agent service started / Upgrade-task cron job started / Refresh-remote-image cron job started / Version-check cron job started / Running OS command:
# 
# ImportError: python import error
# PrepareBundleError: upgrade

# cat result.out | tr -d '\000' | grep -e ABORT -e DA_ASSERT -e dafile -e checkpoint= -e CalledProcessError -e "Caught exception" -e "Handling RpcUpgradeAgent" -e DaExceptionRpcTimeout -e "Killing UpgradeMgr"

# panic, excpetion
grep -e dafile -e ABORT -e DA_ASSERT -e CalledProcessError -e "Caught exception" -e DaExceptionRpcTimeout -e Traceback result.out

# upgrade
# grep -e dafile -e checkpoint= -e "Handling RpcUpgradeAgent" -e "Killing UpgradeMgr" -e "__ERR__ Missing Agent" -e "Current ZK agents count" -e "ZOO_ERROR@" result.out
# Agent(vesx00.datrium.com, host-da:00:37:0f:ec:ad): ZK currentStates updated: RUN [4.0.101.0-31091_a2feaf7_g] DL [9.0.8.0-999999_abcdef] PREP [none] => 
# Kill bug 35996:  Process ['/opt/datrium/bin/da_setup', '--persist_selected_version'](131057) timed out. Kill.
grep -e dafile -e checkpoint= -e "DVX upgrade" -e "preinstall check failed" -e "Failed to switch software" -e "Handling RpcUpgradeAgent" -e "Killing UpgradeMgr" -e "__ERR__ Missing Agent" -e "Current ZK agents count" -e "ZK currentStates updated" -e Kill -e "timed out" result.out

# mac address
grep -e dafile -e "Agent(vesx00.datrium.com" result.out
grep -e dafile -e " Agent(kvmfrontend0, host" result.out

# HA
# received req prepareUpgrade, received req failoverInternal
# resync succeeded, resync fail??
grep -e dafile -e check_upgrade_state -e "fail over" -e "shouldFailover status" -e "received req failoverInternal" -e "received req prepareUpgrade" -e "resync succeeded" -e "succeeded to take over" -e FailoverFinishedEvent -e failup -e CHANNEL_ZOOKEEPER result.out

# ssd issues
# bug 36579
# eventType: ESXHBAStatsFailureEvent, ESXSSDStatsFailureEvent
# vmkernel: PDL error, "unmapped", "event code:", "Path lost", "permanently inaccessible", "Permanent Device Loss", "Device is permanently unavailable", "Failed to write"
# zgrep "ScsiDeviceIO" vmkernel.0.gz | grep -v 0x85 | grep -v 0x1a
# $ pwd
# /cores/bug_36579/a/host/C220M4S03_00_25_b5_00_00_6b/etc/vmware
# $ cat esx.conf | grep -i OptLockReadTimeout
# /adv/VMFS3/OptLockReadTimeout = "5000"
# DeletePath : adapter=vmhba2, channel=0, target=12, lun=0
#
# reset and timeout
# lsi_mr3: megasas_hotplug_work:255: event code: 0x10b
# lsi_mr3: megasas_hotplug_work:255: event code: 0x10c
#
# bug 35912
# 0x2a, 0x28 <= read or write IO errors
# WARNING: ScsiScan: 2011: Could not delete path vmhba2:C0:T12:L0

grep -e dafile -e ESXHBAStatsFailureEvent -e ESXSSDStatsFailureEvent -e "PDL error" -e unmapped -e "Path lost" -e "permanently inaccessible" -e "Permanent Device Loss" -e "Device is permanently unavailable" -e "Failed to write" result.out

# FE
# dlog fe.log.2018-06-05T*.38521.1.gz|grep -i StuckIO
