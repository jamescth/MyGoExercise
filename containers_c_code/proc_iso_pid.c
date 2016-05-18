#define _GNU_SOURCE
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/wait.h>
#include <unistd.h>
#include <errno.h>

static char child_stack[1048576];

static int child_fn() {
		printf("PID: %ld\n", (long)getpid());
		printf("PPID: %ld\n", (long)getppid());
		return 0;
}

int main() {
	pid_t child_pid = clone(child_fn, child_stack+1048576, CLONE_NEWPID | SIGCHLD, NULL);
	if (child_pid == -1) {
			printf("clone() failed, errno = %d\n", errno);
			return errno;
	}

	printf("clone() = %ld\n", (long)child_pid);

	waitpid(child_pid, NULL, 0);
	return 0;
}
