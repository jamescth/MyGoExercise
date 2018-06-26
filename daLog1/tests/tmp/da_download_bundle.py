#!/bin/python

import sys
import time
import hashlib
import subprocess

def terminate_process(process_handle):
    '''
    Terminate process by sending signals to it. See command.py.
    '''
    try:
       # Popen.kill and Popen.terminate don't work on ESX.
       os.kill(process_handle.pid, signal.SIGKILL)
    except:
       pass

def wait_for_process(process_handle, timeout):
    '''
    Check if child process has terminated. Otherwise terminate it
    after timeout seconds.
    '''
    rc = None
    poll_started = time.time()
    while rc is None:
        if time.time() - poll_started > timeout:
            terminate_process(process_handle)
            rc = process_handle.poll()
            break

        time.sleep(0.5)
        rc = process_handle.poll()
    return rc

def return_code_to_message(process_desc, return_code):
    '''
    Convert the process's return code to message.
    '''
    if return_node is None:
        return "%s has not terminated" % process_desc
    elif return_code < 0:
        return "%s terminated by signal %d" % (process_desc, -return_code)
    else:
        return "%s terminated with return code %d" % (process_desc, return_code)

def download_upgrade_bundle(url, bundle_checksum):
    '''
    Download, checksum, and extract the hyperdriver bundle from the DataNode
    without any temporary files.  Requires fewer resources and less cleanup
    when things go wrong.  It may even be faster.

    @param url: location of the hyperdriver bundle on the DataNode.
    @param bundle_checksum: checksum of the hyperdriver bundle.
    '''
    IOSIZE = 64 * 1024
    staging_dir='/da-sys/stage'
    checksum = hashlib.md5()
    extract_bundle = None
    try:
        bundle_download = subprocess.Popen(["wget", "-O-", url],
                                           stdout=subprocess.PIPE, cwd='/',
                                           bufsize=IOSIZE)
        extract_bundle = subprocess.Popen(["tar", "zxf", "-", "-C", staging_dir],
                                          stdin=subprocess.PIPE, cwd='/',
                                          bufsize=IOSIZE)
        download_started = time.time()
        bytes_transferred = 0
        download_stalled = False
        while True:
            bundle_data = bundle_download.stdout.read(IOSIZE)
            if not bundle_data:
                break

            # bootstrap has a timeout of 150 seconds
            if time.time() - download_started > 140:
                download_stalled = True
                break

            bytes_transferred += len(bundle_data)
            checksum.update(bundle_data)
            extract_bundle.stdin.write(bundle_data)

        return_code = wait_for_process(bundle_download, 5)
        if download_stalled:
            raise Exception("download of bundle is too slow")

        if return_code != 0:
            raise Exception(return_code_to_message("download of bundle", return_code))

        print ("Transferred and checksummed %d bytes in %d seconds" %
               (bytes_transferred, time.time() - download_started))
        extract_bundle.stdin.close()
    except:
        # make sure the tar process has exited so we can clean up the ramdisk.
        if extract_bundle is not None:
            try:
                extract_bundle.stdin.close()
            except:
                pass
            wait_for_process(extract_bundle, 5)
        raise
    return_code = wait_for_process(extract_bundle, 5)
    if return_code != 0:
        raise Exception(return_code_to_message("extract of bundle", return_code))

    if bundle_checksum != checksum.hexdigest():
        raise Exception("expected checksum '%s' does not match"
                        " calculated checksum '%s'" % (bundle_checksum,
                        checksum.hexdigest()))

if __name__ == '__main__':
    try:
        download_upgrade_bundle(sys.argv[1].strip(), sys.argv[2].strip())
    except Exception as e:
        # TODO (yang): pass an error to config agent to change the state file.
        print (str(e))
        raise
