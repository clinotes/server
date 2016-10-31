# clinot.es server [![CircleCI](https://circleci.com/gh/clinotes/server.svg?style=svg)](https://circleci.com/gh/clinotes/server)

This is a little side-project to learn `Go` and play with [Heroku](https://heroku.com), [Postmark](https://postmarkapp.com), [Stripe](https://stripe.com) and PostgreSQL.

The service at [clinot.es](https://clinot.es) is a remote note management service with a [CLI client](https://github.com/clinotes/client). Use the hosted service or feel to host it by yourself, it works fine with Heroku and the default PostgreSQL add-on.

## License

Feel free to use the server code, it's released using the [GPLv3 license](https://github.com/clinotes/server/blob/master/LICENSE.md).

## Features

- [x] Create account
- [x] Verify account
- [x] Create authorization token
- [x] Verify authorization token
- [x] Create subscriptions using [Stripe.com](https://stripe.com)
- [ ] Create notes

## Dependecies

- PostgreSQL database
- [Postmark](https://postmarkapp.com) (send emails)
- [Stripe](https://stripe.com) (handle subscriptions)

## Install

There will be an easy guide to install the server someday …
