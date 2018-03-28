import random
import threading
import socket
import datetime
import string
import signal

threads = 10
ip = '192.168.1.152'
port = 4445
BUFFER_SIZE = 1024
lock = threading.Lock()
run = True

def worker():
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.settimeout(20)
    while run:
        try:
            stock = ''
            for _ in range(3):
                stock += random.choice(string.ascii_uppercase)
            start = datetime.datetime.now()
            s.connect((ip, port))
            s.send(stock)
            data = s.recv(BUFFER_SIZE)
            s.close()
            delta = datetime.datetime.now() - start()
            with lock:
                print(delta.milliseconds)
        except KeyboardInterrupt:
            return
        except Exception as e:
            with lock:
                print(e)

def shutdown(signum, frame):
    run = False

if __name__ == '__main__':
    signal.signal(signal.SIGTERM, shutdown)
    signal.signal(signal.SIGINT, shutdown)

    workers = []
    for i in range(threads):
        workers.append(threading.Thread(name='worker', target=worker))

    for worker in workers:
        worker.setDaemon(True)
        worker.start()

    for worker in workers:
        worker.join()