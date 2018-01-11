FROM python:3.5.1

RUN apt-get -q update && apt-get install -y -q \
  sqlite3 --no-install-recommends \

  && apt-get clean && rm -rf /var/lib/apt/lists/*



RUN apt-get update && apt-get install -y \
        binutils \
        libproj-dev \
        gdal-bin \
        libgdal-dev 
ENV LANG C.UTF-8

RUN pip install --upgrade pip virtualenv

RUN virtualenv /venv
ENV VIRTUAL_ENV /venv
ENV PATH /venv/bin:$PATH

RUN mkdir -p /app
WORKDIR /app

ADD requirements.txt /app/requirements.txt
RUN pip install --no-cache-dir -r /app/requirements.txt
RUN pip install django-filter

ADD . /app

CMD python3 manage.py migrate
CMD python3 manage.py populate
CMD python manage.py runserver 0.0.0.0:8000
