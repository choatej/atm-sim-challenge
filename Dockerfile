FROM golang:1.20.4-buster
LABEL authors="jchoate"

WORKDIR /app

COPY atm-sim_linux_amd64 /app/atm-sim
ENTRYPOINT ["/app/atm-sim"]