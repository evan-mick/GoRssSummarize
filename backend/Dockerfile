# https://docs.docker.com/language/golang/build-images/

FROM golang:1.22-alpine3.18
WORKDIR /container
COPY go.mod go.sum ./
COPY ./frontend_data ./frontend_data
COPY ./frontend ./frontend
COPY .env ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /rsssummarize
RUN ["chmod", "+x", "/rsssummarize"]
RUN ls
EXPOSE 9999
EXPOSE 8080

# CMD ["./rsssummarize"]
CMD [ "/rsssummarize" ]