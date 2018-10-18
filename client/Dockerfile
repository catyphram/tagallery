FROM ubuntu

RUN apt-get update
RUN apt-get install -y python3

# EXPOSE  8080

WORKDIR /var/www/
CMD ["python3", "-m", "http.server", "8080"]
