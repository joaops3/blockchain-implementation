FROM golang:1.24

WORKDIR /go/internal
ENV PATH="/go/bin:${PATH}"

COPY . .

CMD ["tail", "-f", "/dev/null"]