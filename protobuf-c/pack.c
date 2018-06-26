#include <stdio.h>
#include <stdlib.h>
#include "amessage.pb-c.h"

int main(int argc, const char * argv[])
{
	AMessage msg = AMESSAGE__INIT;
	void *buf;
	unsigned len;

	if (argc != 2 && argc != 3) {
		// Allow one or two integers
		fprintf(stderr, "usage: amessage a [b]\n");
		exit(1);
	}

	msg.a = atoi(argv[1]);
	if (argc == 3) {
		msg.has_b = 1;
		msg.b = atoi(argv[2]);
	}
	len = amessage__get_packed_size(&msg);

	buf = malloc(len);
	amessage__pack(&msg, buf);

	fprintf(stderr, "Writing %d serialized bytes\n", len);
	fwrite(buf, len, 1, stdout);

	free(buf);
	// exit(0);
}
