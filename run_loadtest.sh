#!/bin/bash
iterations=$1

# This script just runs the load test $iteration times and saves the metrics to a text file
for i in `eval echo {1..$iterations}`
do
	/home/ubuntu/she-chain/build/loadtest
	date | tee -a loadtest_results.txt
	sleep 5
	python3 /home/ubuntu/she-chain/loadtest/scripts/metrics.py | tee -a loadtest_results.txt
done

