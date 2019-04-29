FROM golang 

EXPOSE 8081

RUN ["/usr/local/go/bin/go","get", "-u", "github.com/gorilla/mux"]
RUN ["/usr/local/go/bin/go","get", "-u", "github.com/lib/pq"]

ADD ./main.go /usr/local/bin/main.go

CMD ["go", "run", "/usr/local/bin/main.go"]

