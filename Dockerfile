FROM scratch
LABEL maintainer="tyr@pettingzoo.co"

EXPOSE 8080

ENV DB_ENGINE  postgresql://postgres:5432/api
ENV REDIS_ADDR localhost:6379

ADD puphaus-api /
CMD ["/puphaus-api"]