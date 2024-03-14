FROM node:lts AS frontend-build
WORKDIR /app
COPY . .
RUN npm i
RUN npm run build

FROM golang:1.22-alpine as backend-builder

RUN mkdir -p /otennie
WORKDIR /otennie

COPY backend/ .

RUN go mod download

RUN CGO_ENABLED=0 go build -a -o bin/server cmd/server/*

FROM alpine:3.19

RUN addgroup -S app \
  && adduser -S -G app app 

WORKDIR /home/app

COPY --from=backend-builder /otennie/bin/server .
COPY --from=frontend-build /app/dist .
RUN chown -R app:app ./

USER app

CMD ["./server"]
