FROM golang:1.20

ARG db_dialect
ARG db_host
ARG db_port
ARG db_user
ARG db_password
ARG db_name

RUN apt update && apt install -y wkhtmltopdf

RUN mkdir -p /go/src/app

WORKDIR /go/src/app

COPY . .

RUN go build -o goreports

RUN chmod +x goreports

RUN ./goreports init --db-dialect=$db_dialect --db-host=$db_host --db-port=$db_port --db-username=$db_user --db-password=$db_password --db-name=$db_name

EXPOSE 3200

CMD ["./goreports", "start"]