# -- BUILD STAGE --
FROM golang:1.21.6-alpine3.18 as BUILD 

WORKDIR /build

#RUN go install github.com/hashicorp/terraform-exec/tfexec@latest

COPY . .

RUN go build -o ihp-terraform

## -- RUNTIME STAGE --
FROM alpine:3.18.4 as RUNTIME

WORKDIR /app

# Copy build with permissions
COPY --from=BUILD /build/ihp-terraform /app/ihp-terraform

# Copy TerraDocker
COPY ./TerraDocker /app/TerraDocker

# Add volume for Terraform files
VOLUME /app/TerraDocker

CMD ["/app/ihp-terraform"]