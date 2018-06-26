#!/bin/bash
#./daLog -lpath /cores/bug_30350
#./daLog -lpath /cores/bug_29487

#./daLog -is_http 1 
# cat result.out | cnul | grep -e ABORT -e DA_ASSERT
cat result.out | tr -d '\000' | grep -e ABORT -e DA_ASSERT -e Traceback -e "Exception:"
