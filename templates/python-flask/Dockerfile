FROM jetpackio/devbox:latest

WORKDIR /code
USER root:root
RUN mkdir -p /code && chown ${DEVBOX_USER}:${DEVBOX_USER} /code
USER ${DEVBOX_USER}:${DEVBOX_USER}

COPY --chown=${DEVBOX_USER}:${DEVBOX_USER} . .

RUN devbox install

SHELL ["/bin/bash", "-c"]
RUN devbox run install

RUN echo 'No build script found, skipping'

CMD ["devbox", "run", "start"]
