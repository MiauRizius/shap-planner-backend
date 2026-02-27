# Wir sagen Docker: Der Builder soll IMMER auf der Architektur deines PCs laufen (schnell!)
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

RUN apk add --no-cache git
WORKDIR /app

# Cache für Module nutzen
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# buildx übergibt diese Variablen automatisch
ARG TARGETOS
ARG TARGETARCH

# Hier passiert die Magie: Go kompiliert NATIV für das Ziel (Cross-Compilation)
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags "-s -w" -o shap-planner-backend .

# Final Stage (bleibt gleich klein)
FROM scratch
COPY --from=builder /app/shap-planner-backend /shap-planner-backend
ENTRYPOINT ["/shap-planner-backend"]