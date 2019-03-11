FROM ubuntu:18.04

RUN apt-get -y update

ENV PGVER 10
ENV GOVER 1.10
RUN apt-get install -y postgresql-$PGVER

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

USER root

RUN apt install -y golang-$GOVER git

ENV GOROOT /usr/lib/go-$GOVER
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR $GOPATH/src/2019_1_Auteam
ADD . $GOPATH/src/2019_1_Auteam

RUN go install ./
RUN go install ./sessions_app/session_server/

EXPOSE 8080

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    /etc/init.d/postgresql stop

EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

CMD service postgresql start && ./session_app/session_server/session_server && 2019_1_Auteam