# https://docs.docker.com/language/golang/build-images/

FROM golang:1.24
WORKDIR /container
COPY go.mod go.sum ./
COPY points.json ./
COPY ./frontend_data ./frontend_data
COPY ./frontend ./frontend
COPY .env ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /rsssummarize
RUN ["chmod", "+x", "/rsssummarize"]
RUN ls
EXPOSE 8080

# CMD ["./rsssummarize"]
CMD [ "/rsssummarize" ]
