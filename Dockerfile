FROM golang

WORKDIR /GolangWebApp

COPY go.mod go.sum ./
RUN go mod download && go mod verify

EXPOSE 4040

COPY . .
RUN make fbuild

CMD ["./apiserver", "config-path"]