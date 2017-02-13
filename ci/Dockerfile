FROM vendasta/mscli:1.16.0
COPY ./serviceaccount /var/run/secrets/kubernetes.io/serviceaccount/
COPY hal9000.json /etc/hal9000/hal9000.json
COPY . /go/src/github.com/vendasta/event-store
