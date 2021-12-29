# build binary
FROM golang:1.17-alpine3.15 AS build
RUN apk add git
WORKDIR /go/mod/github.com/dieqnt/skeleton
COPY . /go/mod/github.com/dieqnt/skeleton
RUN go mod download
RUN CGO_ENABLED=0 go build -o /out/templatemicroservice github.com/dieqnt/skeleton/cmd/templatemicroserviced

COPY db/elasticsearch/mappings/visenze_index_v0-mapping.json /out
COPY db/elasticsearch/mappings/visenze_index_v0-settings.json /out

# copy to alpine image
FROM alpine:3.15
WORKDIR /app
COPY --from=build /out/visenze_index_v0-mapping.json /app/db/elasticsearch/mappings/visenze_index_v0-mapping.json
COPY --from=build /out/visenze_index_v0-settings.json /app/db/elasticsearch/mappings/visenze_index_v0-settings.json
COPY --from=build /out/templatemicroservice /app
CMD ["/app/templatemicroservice"]
# docker build -t templatemicroservice .
