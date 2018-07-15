FROM alpine:3.4

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

EXPOSE 8000

ADD gameapi /usr/local/bin/gameapi

CMD ["gameapi"]
