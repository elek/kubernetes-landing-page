FROM            golang:1.10.0-alpine3.7
RUN             apk add --update --no-cache git make curl
RUN				 curl https://glide.sh/get | sh


                    

WORKDIR         /go/src/kubernetes-landing-page
COPY            . .
RUN             make


FROM            golang:1.10.0-alpine3.7
COPY            --from=0 /go/bin/kubernetes-landing-page /go/bin/kubernetes-landing-page
EXPOSE          8080
ENV             GIN_MODE=release

ENTRYPOINT      [ "/go/bin/kubernetes-landing-page" ]
