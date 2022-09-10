FROM golang
ADD . /usr/src/soap
WORKDIR /usr/src/soap
RUN go build -v -o /usr/local/bin/soap ./
CMD ["soap"]