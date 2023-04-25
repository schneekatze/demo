FROM golang:1.19-alpine

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

WORKDIR /opt/chlng

COPY . .

RUN go build .
USER nobody:nobody
EXPOSE 8000

CMD ["/opt/chlng/challenge", "server", "serve"]
