
from threading import Thread
from threading import Lock

i = 0
mtx = Lock()

def adder():
	global i

	for x in range(0,1000000):
		mtx.acquire()
		i+=1
		mtx.release()

def subber():
	global i
	
	for x in range(0,1000000):
		mtx.acquire()
		i-=1
		mtx.release()

def main():
	adder_thr = Thread(target = adder)
	subber_thr = Thread(target = subber)

	adder_thr.start()
	subber_thr.start()

	adder_thr.join()
	subber_thr.join()

	print("Done: " + str(i))

main()
