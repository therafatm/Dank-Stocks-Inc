import numpy as np
import matplotlib.pyplot as plt
import sys

if len(sys.argv) < 2:
	print "Please provide data filepath"

data = np.genfromtxt(sys.argv[1])
# data = data[data>2000]
plt.hist(data)
plt.show()