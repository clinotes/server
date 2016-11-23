# clinot.es server [![CircleCI](https://circleci.com/gh/clinotes/server.svg?style=svg)](https://circleci.com/gh/clinotes/server)

This is a side-project in `Go` to play with [Heroku](https://heroku.com), [Postmark](https://postmarkapp.com), [Stripe](https://stripe.com) and PostgreSQL.

The server at [clinot.es](https://clinot.es) is a remote note management service with a [command line client](https://github.com/clinotes/client). Use the hosted service or feel to host it by yourself …

## Features

- [x] Create account
- [x] Verify account
- [x] Create access token
- [x] Verify access token
- [x] Create subscriptions (draft)
- [x] Create notes
- [x] List notes

## Dependecies

- PostgreSQL database
- [Postmark](https://postmarkapp.com) (send emails, **required**)
- [Stripe](https://stripe.com) (handle subscriptions, **draft**)

## Setup

The API server works fine on a free Heroku dyno using the free PostgreSQL add-on. The free [Postmark](https://postmarkapp.com) accounts comes with 25.000 free emails, so this should last for a couple of accounts ;)

### Postmark

[Postmark](https://postmarkapp.com) is used for sending emails to new users. You need to create three templates in your Postmark account and configure the template IDs in your environment variables. You will find the three HTML and plaintext templates inside the `templates/` folder:

* [Welcome](/templates/welcome)
* [Confirmation](/templates/confirmation)
* [Access Token](/templates/token)

Make sure to validate your sender address in Postmark as well!

### Application

```bash
$ > git clone git@github.com:clinotes/server.git
$ > cd server
$ > heroku create
$ > heroku addons:create heroku-postgresql:hobby-dev
$ > git push heroku master
$ > heroku info

…
Web URL: https://exmaple-url-12345.herokuapp.com/
```

### Environment

```bash
$ > heroku config:set MAX_DB_CONNECTIONS=5
$ > heroku config:set POSTMARK_API_KEY=API_KEY
$ > heroku config:set POSTMARK_TEMPLATE_WELCOME=TEMPLATE_ID
$ > heroku config:set POSTMARK_TEMPLATE_CONFIRMATION=TEMPLATE_ID
$ > heroku config:set POSTMARK_TEMPLATE_TOKEN=TEMPLATE_ID
$ > heroku config:set POSTMARK_FROM=mail@clinot.es
$ > heroku config:set POSTMARK_REPLY_TO='"CLI Notes" <mail@clinot.es>'
```

### Client

```
$ > brew tap clinotes/cn
$ > brew install cn
$ > echo "CLINOTES_API_HOSTNAME: https://exmaple-url-12345.herokuapp.com/" > ~/.clinotes.yaml
$ > cn version

client: v0.2.1
server: v0.0.5 (supports client >= v0.1.0)
```

## License

Feel free to use the server code, it's released using the [GPLv3 license](https://github.com/clinotes/server/blob/master/LICENSE.md).

## Contributors

- [Sebastian Müller](https://sbstjn.com)
