
from threading import Thread

i = 0

def adder():
	global i

	for x in range(0,1000000):
		i+=1

def subber():
	global i
	
	for x in range(0,1000000):
		i-=1

def main():
	adder_thr = Thread(target = adder)
	subber_thr = Thread(target = subber)

	adder_thr.start()
	subber_thr.start()

	adder_thr.join()
	subber_thr.join()

	print("Done: " + str(i))

main()
