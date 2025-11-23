FROM golang
WORKDIR /app
COPY . .
RUN go mod download
RUN go build
CMD [ "./systm-go-social" ]