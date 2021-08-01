FROM postgres
WORKDIR /docker-entrypoint-initdb.d/
COPY init.sql ./