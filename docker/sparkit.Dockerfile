FROM alpine:latest

ENV EXECUTABLE=sparkit

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app

WORKDIR /app

COPY --from=sparkit-builder:latest /application/${EXECUTABLE} /app

CMD /app/${EXECUTABLE} 
