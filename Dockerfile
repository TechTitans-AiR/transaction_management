FROM golang:latest

# Postavi radni direktorij
WORKDIR /go/src/app

# Kopiraj POM i sve potrebne datoteke za preuzimanje dependencija
COPY . .

# Postavi varijable okoline koje će biti dostupne tijekom izvođenja Docker slike
ARG MONGO_URI

# Stvori .env datoteku unutar Docker kontejnera
RUN echo "MONGO_URI=${MONGO_URI}" > .env

# Pokreni aplikaciju
CMD ["go", "run", "main.go"]