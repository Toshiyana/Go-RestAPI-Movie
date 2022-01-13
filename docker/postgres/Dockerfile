FROM postgres:14-alpine
ENV LANG ja_JP.utf8

# Initialize database by copying sql files
COPY ./docker/postgres/init/*.sql /docker-entrypoint-initdb.d/