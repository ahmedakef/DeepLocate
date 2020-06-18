FROM golang:latest


RUN apt-get update && apt-get install -y \
    python3 python3-pip
RUN ln -sf /usr/bin/python3 /usr/bin/python

RUN pip3 install --upgrade pip --user

WORKDIR /app

COPY ImageCaptioningAndKeyWordExtraction/requirements.txt /app/ImageCaptioningAndKeyWordExtraction/requirements.txt
RUN pip3 install -r ImageCaptioningAndKeyWordExtraction/requirements.txt
RUN python3 -m spacy download en_core_web_sm

COPY ./ /app

WORKDIR /app/dlocate


RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

# RUN go build

ENTRYPOINT CompileDaemon
# ENTRYPOINT ./dlocate -o index -d /home/
# CMD [""]