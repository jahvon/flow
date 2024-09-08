FROM golang:1.23.1-bookworm

WORKDIR /workspaces
ENV DISABLE_FLOW_INTERACTIVE="true"

# TODO: replace with examples repo
ENV WORKSPACE="flow"
ENV REPO="https://github.com/jahvon/flow.git"
ENV BRANCH=""

COPY . flow
WORKDIR /workspaces/flow
RUN go build -o /usr/local/bin/flow main.go && chmod +x /usr/local/bin/flow && rm -rf /workspaces/flow

WORKDIR /workspaces
RUN if [ -z "$BRANCH" ]; then git clone $REPO .; else git clone -b $BRANCH $REPO .; fi
RUN flow init workspace $WORKSPACE . --set

ENTRYPOINT ["flow"]
CMD ["get", "workspace"]