# build binary
ARG DOCKER_PROXY
FROM ${DOCKER_PROXY}/golang:1.20-alpine3.17 AS build
RUN apk add git
RUN apk --update add tzdata git
RUN cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime
RUN echo "Europe/Moscow" > /etc/timezone
RUN date
RUN apk del tzdata

ARG GITLAB_USER
ARG GITLAB_TOKEN
ENV GO111MODULE=on
ENV GOPRIVATE="gitlab.goodsteam.tech/*"
RUN echo -e "machine gitlab.goodsteam.tech\nlogin ${GITLAB_USER}\npassword ${GITLAB_TOKEN}" > /root/.netrc

WORKDIR /go/mod/github.com/pocoz/skeleton
COPY . /go/mod/github.com/pocoz/skeleton
RUN go mod download
RUN CGO_ENABLED=0 go build -o /out/skeleton github.com/pocoz/skeleton/cmd/skeletond

COPY db/elasticsearch/mappings/visenze_index_v0-mapping.json /out
COPY db/elasticsearch/mappings/visenze_index_v0-settings.json /out

# copy to alpine image
FROM ${DOCKER_PROXY}/alpine:3.17
WORKDIR /app
COPY --from=build /out/visenze_index_v0-mapping.json /app/db/elasticsearch/mappings/visenze_index_v0-mapping.json
COPY --from=build /out/visenze_index_v0-settings.json /app/db/elasticsearch/mappings/visenze_index_v0-settings.json
COPY --from=build /out/skeleton /app
CMD ["/app/skeleton"]
# docker build -t skeleton .
