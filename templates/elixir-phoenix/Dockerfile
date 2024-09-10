FROM ubuntu:latest

# If using a different user, then may need to create the user below.
ARG DEVBOX_USER=root

USER root:root

RUN apt-get update
# This issue has some tips for installing Elixir with Ubuntu
# https://github.com/phoenixframework/phoenix/issues/5552
# erlang-dev for SSL libraries
# erlang-xmerl for missing xmerl.app error
RUN apt-get install -y elixir git erlang-dev erlang-xmerl

WORKDIR /code

USER ${DEVBOX_USER}:${DEVBOX_USER}
RUN mkdir -p /code && chown ${DEVBOX_USER}:${DEVBOX_USER} /code

COPY --chown=${DEVBOX_USER}:${DEVBOX_USER} . .

# Install hex (package manager) and rebar (erlang build tool) locally
RUN mix local.hex --force
RUN mix local.rebar --force

RUN mix setup

CMD mix phx.server
