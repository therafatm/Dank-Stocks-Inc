import time
import psycopg2
from http.server import BaseHTTPRequestHandler, HTTPServer
import json
import threading
from os import environ
import queue

masters = {}
q = queue.Queue()

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

def wait_for_connect(host, port, passwd, user, db, retrys=20):
    for i in range(retrys):
        try:
            conn = psycopg2.connect("dbname=%s user=%s host=%s port=%s password=%s" % (db, user, host, port, passwd))
            conn.autocommit = True
            print("connected " + host + ":" + port)
            return conn
        except:
            print("retrying " + host + ":" + port)
            time.sleep(1)

    print("Timed out connecting to " + host + ":" + port)
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
        conn = wait_for_connect(data['master_host'], data['master_port'], data['password'], data['user'], data['db'])
        worker_conn = wait_for_connect(data['worker_host'], data['worker_port'], data['password'], data['user'], 'postgres')
        
        if conn and worker_conn:
            master = data['master_host'] + ':' + data['master_port']

            cur = conn.cursor()
            cur.execute("""SELECT master_add_node(%(worker_host)s, %(worker_port)s);""", data)
            print("adding " + data['worker_host'] + ":" + data['worker_port'] + " to " + master)
                
            if master not in masters:
                cur.execute("""SELECT distribute();""")
                masters[master] = 0

            masters[master] += 1

        else:
            print("Timed out.")


if __name__ == "__main__":
    worker_thread = threading.Thread(name='worker', target=process)
    server_thread = threading.Thread(name='server', target=run)
    worker_thread.start()
    server_thread.start()

    worker_thread.join()
    server_thread.join()