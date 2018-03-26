FROM python:3
RUN pip install psycopg2-binary
ADD manager /manager
ENTRYPOINT python -u /manager/manager.py