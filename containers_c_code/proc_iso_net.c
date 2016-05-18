#define _GNU_SOURCE
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/wait.h>
#include <unistd.h>
#include <errno.h>

static char child_stack[1048576];

static int child_fn() {
		printf("New 'net' Namespece:\n");
		system("ip link");
		printf("\n\n");
		return 0;
}

int main() {
	printf("Original 'net' Namespace:\n");
	system("ip link");
	printf("\n\n");

	pid_t child_pid = clone(child_fn, child_stack+1048576,
				   	CLONE_NEWPID | CLONE_NEWNET | SIGCHLD, NULL);
	if (child_pid == -1) {
			printf("clone() failed, errno = %d\n", errno);
			return errno;
	}

	waitpid(child_pid, NULL, 0);
	return 0;
}
