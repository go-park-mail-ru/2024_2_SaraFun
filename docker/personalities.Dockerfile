FROM alpine:latest

ENV EXECUTABLE=personalities

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app

WORKDIR /app

COPY --from=sparkit-personalities-builder:latest --chmod=755 /application/${EXECUTABLE} /app

CMD /app/${EXECUTABLE}