# Python/Flask Example Project

[![Deploy on Jetify](https://www.jetify.com/img/devbox/deploy-on-jetify-galaxy.svg)](https://cloud.jetify.com/team/new/deploy?repo=github.com/jetify-examples/python-flask)

[Flask](https://flask.palletsprojects.com/en/3.0.x/) is a microframework for Python web services. This repo is an example project for use with [Jetify Cloud](https://www.jetify.com/cloud)

## How to Use

* Create an account on [Jetify Cloud](https://cloud.jetify.com)
* Fork or clone this repo to your own account
* In the Jetify Dashboard, create a new project, and select the cloned repo from your account

## About this project

Note that this project creates a `venv` automatically, and activates it before running any commands or shells.

### Packages Installed

* Python@3.11
* Gunicorn
* Flask

### Scripts

Install: `pip install -r requirements.txt`

Start: `gunicorn -w 2 -b 0.0.0.0:8080 main:app`

