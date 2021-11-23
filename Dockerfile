FROM alpine:3.13.7

ENV USER_UID=1001 \
    USER_NAME=admin-console \
    HOME=/home/admin-console \
    PATH="/go/bin:${PATH}"

RUN addgroup --gid ${USER_UID} ${USER_NAME} \
    && adduser --disabled-password --uid ${USER_UID} --ingroup ${USER_NAME} --home ${HOME} ${USER_NAME}

WORKDIR /go/bin

COPY go-binary .
COPY static static
COPY views views
COPY conf conf
COPY db db

USER ${USER_UID}

ENTRYPOINT ["go-binary"]
