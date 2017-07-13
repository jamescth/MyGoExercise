#!/bin/sh
# Identify allocations
# https://methane.github.io/2015/02/reduce-allocation-in-go-code/
GODEBUG=allocfreetrace=1 ./ptest.test -test.bench=BenchmarkSecond -test.benchtime=10ms 2>trace.log

# inline
# http://www.golangbootcamp.com/book/tricks_and_tips
# https://dave.cheney.net/2014/06/07/five-things-that-make-go-fast
go build -gcflags=-m main.go
