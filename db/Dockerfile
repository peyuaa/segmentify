FROM postgres:15.4

ENV POSTGRES_HOST_AUTH_METHOD trust

COPY init.sql /docker-entrypoint-initdb.d/

EXPOSE 5432