FROM golang:1.11

RUN echo "[Installing dependencies]"
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN cd /go/src
WORKDIR /go/src
RUN git clone -b RouterDev https://github.com/go-park-mail-ru/2019_1_Auteam.git
WORKDIR /go/src/2019_1_Auteam
RUN dep ensure

RUN echo "[Building golang application]"
RUN go build ./apps/sessions_app/sessions_server
COPY init.sh /go/init.sh
CMD /go/init.sh

EXPOSE 8081
