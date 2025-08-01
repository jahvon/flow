FROM golang:1.24.5-bookworm

ENV DISABLE_FLOW_INTERACTIVE="true"

# TODO: replace with examples repo
ENV WORKSPACE="flow"
ENV REPO="https://github.com/flowexec/flow.git"
ENV BRANCH=""

WORKDIR /workspaces
COPY flow /usr/bin/flow

RUN if [ -z "$BRANCH" ]; then git clone $REPO .; else git clone -b $BRANCH $REPO .; fi
RUN flow workspace create $WORKSPACE . --set

ENTRYPOINT ["flow"]
CMD ["--version"]