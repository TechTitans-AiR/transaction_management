FROM golang:latest

# Postavi radni direktorij
WORKDIR /go/src/app

# Kopiraj POM i sve potrebne datoteke za preuzimanje dependencija
COPY . .

# Postavi varijable okoline koje će biti dostupne tijekom izvođenja Docker slike
ARG MONGO_URI
ARG PORT

# Stvori .env datoteku unutar Docker kontejnera
RUN echo "MONGO_URI=${MONGO_URI}"
RUN echo "PORT=${PORT}"

# Pokreni aplikaciju
CMD ["go", "run", "main.go"]