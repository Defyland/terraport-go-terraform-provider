FROM golang:1.25 AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/terraform-provider-terraport ./cmd/terraform-provider-terraport

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/terraform-provider-terraport /terraform-provider-terraport
USER nonroot:nonroot
ENTRYPOINT ["/terraform-provider-terraport"]
