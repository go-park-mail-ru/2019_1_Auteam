FROM postgres:11
COPY pg_hba.conf /pg_hba.conf
COPY postgresql.conf /postgresql.conf
COPY init.sql /docker-entrypoint-initdb.d/
EXPOSE 5432
