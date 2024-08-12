FROM postgres:13

COPY ./schemas/*.sql /docker-entrypoint-initdb.d/