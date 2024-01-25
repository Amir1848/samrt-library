FROM postgres:alpine

ENV POSTGRES_PASSWORD=1
ENV POSTGRES_USER=postgres
ENV POSTGRES_DB=sepidar_library

VOLUME [ "/db-volume/:/var/lib/postgresql/data" ]