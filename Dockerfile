FROM golang:1.9.0-stretch
RUN apt update && apt install -yqq git
WORKDIR /go/src/github.com/elastic/beats
RUN git clone https://github.com/elastic/beats.git && cd beats && git checkout 5.5
WORKDIR /go/src/github.com/PPACI/krakenbeat/
COPY . ./
RUN make
RUN chmod go-w krakenbeat.yml
CMD ./krakenbeat -e
