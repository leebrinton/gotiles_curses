# Build containter
FROM golang:1.26-trixie AS builder
RUN apt-get update && \
    apt-get install -y git ncurses-dev  

WORKDIR /go/src
RUN git clone https://github.com/leebrinton/gotileslib.git; \
    git clone https://github.com/leebrinton/gotiles_curses.git;

WORKDIR /go/src/gotiles_curses
RUN go mod tidy; \
    go build

FROM debian:trixie
RUN apt-get update && \
    apt-get install -y libncurses6
COPY --from=builder /go/src/gotiles_curses/gotiles_curses /usr/local/bin/tiles
CMD ["/usr/local/bin/tiles"]
