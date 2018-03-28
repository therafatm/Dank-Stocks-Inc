import time
import psycopg2
from http.server import BaseHTTPRequestHandler, HTTPServer
import json
import threading
from os import environ
import queue
import socket

masters = {}
q = queue.Queue()
healthq = queue.Queue()
lock = threading.Lock()


#SELECT master_get_active_worker_nodes();
# curl -d '{"master_host":"logdb", "master_port":"5432", "worker_host":"log_worker", "worker_port": "5432", "db":"logs", "user":"postgres", "password":"postgres"}' -H "Content-Type: application/json/" -X POST http://localhost:3000


class RequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        # Send response status code
        self.send_response(200)

        # Send headers
        self.send_header('Content-type','text/html')
        self.end_headers()

        # Send message back to client
        message = "I'm Alive!\n"
        # Write content as utf-8 data
        self.wfile.write(bytes(message, "utf8"))
        return

    def do_POST(self):
        self.data_string = self.rfile.read(int(self.headers['Content-Length']))
        self.send_response(200)
        self.end_headers()
        print(self.data_string)
        data = json.loads(self.data_string.decode("utf8"))
        print(data)
        q.put(data)
        self.wfile.write(bytes("Sure thing buddy!\n", "utf8"))

def wait_for_connect(host, port, passwd, user, db, retrys=60):
    for i in range(retrys):
        try:
            conn = psycopg2.connect("dbname=%s user=%s host=%s port=%s password=%s" % (db, user, host, port, passwd))
            conn.autocommit = True
            return conn
        except:
            time.sleep(1)
    return None


def run(server_class=HTTPServer, handler_class=RequestHandler, port=3000):
    port = environ.get('MANAGER_PORT', port)
    server_address = ('', port)
    httpd = server_class(server_address, handler_class)
    print('Starting httpd...')
    httpd.serve_forever()

def process():
    print ("Worker started")
    while True:
        data = q.get()
        try:
            print (data['worker_host'])
            data['worker_host'] = socket.gethostbyaddr(data['worker_host'])[0]
            print (data['master_host'])
            #data['master_host'] = socket.gethostbyaddr(data['master_host'])[0]
            conn = wait_for_connect(data['master_host'], data['master_port'], data['password'], data['user'], data['db'])
            worker_conn = wait_for_connect(data['worker_host'], data['worker_port'], data['password'], data['user'], 'postgres')

            if conn and worker_conn:
                master = data['master_host'] + ':' + data['master_port']
                if master not in masters:
                    with lock:
                        masters[master] = 0

                cur = conn.cursor()
                cur.execute("""SELECT nodename, nodeport from pg_dist_node where nodename=%(worker_host)s and nodeport=%(worker_port)s;""", data)
                row = cur.fetchone();
                if row:
                    print("Adding back startup {}:{} to active for {}".format(data['worker_host'], data['worker_port'], data['master_host']))
                    cur.execute("""SELECT master_activate_node(%(worker_host)s, %(worker_port)s);""", data)
                    with lock:
                        masters[master] += 1

                else:
                    cur.execute("""SELECT master_add_node(%(worker_host)s, %(worker_port)s);""", data)
                    print("Adding new {}:{} to {}".format(data['worker_host'], data['worker_port'], data['master_host']))
                    with lock:
                        masters[master] += 1
                    try:
                     cur.execute("""SELECT distribute();""")
                    except psycopg2.Error as e:
                        print(e)

                healthq.put(data)

            else:
                print("Timed out.")

        except Exception as e:
            print(e)
            print("Retrying in 5 seconds")
            time.sleep(5)
            q.put(data)

def healthcheck():
    while True:
        data = healthq.get()
        conn = wait_for_connect(data['master_host'], data['master_port'], data['password'], data['user'], data['db'], retrys=2)
        if conn:
            master = data['master_host'] + ':' + data['master_port']
            cur = conn.cursor()
            cur.execute("SELECT master_get_active_worker_nodes();")
            rows = cur.fetchall()
            for row in rows:
                strpair = row[0].replace('(', '').replace(')', '').split(',')
                worker_host = strpair[0]
                worker_port = strpair[1]
                wconn = wait_for_connect(worker_host, worker_port, data['password'], data['user'], data['db'], retrys=2)
                if not wconn:
                    with lock:
                        worker_dict = {'host': worker_host, 'port': worker_port}
                        print("Disabling node {}:{} for {}".format(worker_host, worker_port , data['master_host']))
                        cur.execute("""SELECT master_disable_node(%(host)s, %(port)s)""", worker_dict)
                        masters[master] -= 1

            #check if inactive is alive
            cur.execute("""SELECT nodename, nodeport from pg_dist_node where isactive='f';""")
            rows = cur.fetchall()
            for row in rows:
                worker_host = row[0]
                worker_port = str(row[1])
                wconn = wait_for_connect(worker_host, worker_port, data['password'], data['user'], data['db'], retrys=2)
                if wconn:
                    with lock:
                        worker_dict = {'host': worker_host, 'port': worker_port}
                        print("Adding back {}:{} to active {}".format(worker_host, worker_port, data['master_host']))
                        cur.execute("""select master_activate_node(%(host)s, %(port)s);""", worker_dict)
                        masters[master] -= 1


            healthq.put(data)
        else:
            print("Master down %s:%s".format(data['master_host'], data['master_port']))

    sleep(5)

if __name__ == "__main__":
    n = 5
    workers = []
    healthworkers  = []
    for i in range(n):
        workers.append(threading.Thread(name='worker', target=process))
    for i in range(n):
        healthworkers.append(threading.Thread(name='healthcheck', target=healthcheck))
    server_thread = threading.Thread(name='server', target=run)

    for worker in workers:
        worker.start()

    for worker in healthworkers:
        worker.start()

    server_thread.start()

    for worker in workers:
        worker.join()
    for health in healthworkers:
        worker.join()
    server_thread.join()
