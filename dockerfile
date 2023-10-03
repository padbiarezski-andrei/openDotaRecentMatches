FROM golang:1.21
WORKDIR /usr/src/app
COPY go.mod main.go handlers.go matchesInfo.go playerInfo.go shared.go matches.html players.html players.json matches.html secret.txt ./
RUN go build -v 
CMD ["./openDotaResentmatches"]
EXPOSE 8080