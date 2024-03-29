# -- BUILD STAGE --
FROM golang:1.21.5-alpine3.18 as BUILD 

WORKDIR /build

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

# Build swagger
RUN swag init

RUN go build -o ihp-backend

## -- RUNTIME STAGE --
FROM alpine:3.18.4 as RUNTIME

WORKDIR /app

ARG USER=docker
ARG UID=5432
ARG GID=5433

# Create user for execution

#User group has same name as user
RUN addgroup -g $GID $USER 


RUN adduser \
    --disabled-password \
    --gecos "" \
    --ingroup "$USER" \
    --no-create-home \
	--uid "$UID" \
    "$USER"

# Copy build with permissions
COPY --from=BUILD --chown=$USER:$USER /build/ihp-backend /app/ihp-backend 

# Ensure that backend can be run
RUN chmod +x /app/ihp-backend

USER $USER 

ENV GIN_MODE=release
ENV PORT=8081

CMD ["/app/ihp-backend"]
