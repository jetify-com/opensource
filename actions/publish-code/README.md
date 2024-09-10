# Publish Code (Github Action)

### Publish code between repositories while preserving history

![License: Apache 2.0](https://img.shields.io/github/license/jetify-com/action-move-code)

## What is it?

A Github action that makes it easy to publish code from one git repository to another
while preserving history. One of the primary use cases is to mirror code from a directory
in a monorepo, and publish it as a standalone repository.

This action was originally built by [jetify](https://www.jetify.com). We
do most of our development in an [opensource monorepo](https://github.com/jetify-com/opensource),
and often publish some projects as separate repositories.

In fact, this very repository is published using this Github Action.

## Related Work

-   [Copybara](https://github.com/google/copybara): A tool written by Google to move
    code between repositories. It is written in Java, and thus a bit heavyweight for
    direct use. That said, a future version of this action could package Copybara in
    an easier to use way.
