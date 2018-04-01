import sys
import numpy as np
import matplotlib.pyplot as plt
import random

if len(sys.argv) < 2:
	print "Please provide data filepath"

n = 10000
distribution = np.genfromtxt(sys.argv[1])
print len(distribution)

def connect():
	return distribution[random.randint(0, len(distribution) - 1)]

def fetch(max_retry, base_timeout, add_timeout):
	connections = 0
	total_time = 0

	for attempt in range(0, int(max_retry)):
		response = connect()
		connections += 1

		cur_timeout = base_timeout + (attempt * add_timeout)
		if response < cur_timeout:
			total_time += response
			return connections, total_time, False

		else:
			total_time += cur_timeout

	return connections, total_time, True

def model(max_retry=5, base_timeout=1000., add_timeout=5000.):
	total_aborts = 0
	total_connections = 0
	times = []
	for i in range(n):
		connections, time, aborted = fetch(max_retry, base_timeout, add_timeout)
		total_connections += connections
		times.append(time)
		if aborted:
			total_aborts += 1
	return np.average(times), float(total_aborts) / n, float(total_connections) / n


retry_range = range(0, 10)
base_range = np.arange(0.0, 5000., 50.)
add_range = [0.0] # np.arange(0.0, 10000., 250.)

print retry_range
print base_range
print add_range

results = []

count = 0
total = len(retry_range) * len(base_range) * len(add_range)

for retry in retry_range:
	for base in base_range:
		for add in add_range:
			avg, aborts, conns = model(retry, base, add)
			result = {
				'average_response_time': avg,
				'aborts_per_quote': aborts,
				'conns_per_quote': conns,
				'base_timeout': base,
				'add_timeout': add,
				'max_retry': retry
			}
			results.append(result)
			count += 1
			print str(count) + " / " + str(total)


results.sort(key=lambda res: res['average_response_time'])
results = filter(lambda res: res['aborts_per_quote'] < 0.05, results)

for i in results[0:10]:
	print i