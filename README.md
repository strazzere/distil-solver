# distil-solver

Solve distil networks "anti-bot" detection for usage in bots.

## How to

_**Change the default keys for admin, users, and banned keys, unless you
want to just let anyone use the service.**_ Also change the address `127.0.0.1:3000` in the `src/static/index.html` page if you want pretty statistics on a page instead of just using the API. Then look in the `docs/` directory.

Normally you will want to setup something like [Praxis](https://github.com/strazzere/praxis) for usage of keeping your server ip address clean and ban free when interacting with potentially hostile APIs.

## What is this?

A few years ago I needed to scrape a handful of sites for threat intel reasons and designed a system to plug into some of my other automation frameworks. Later, a client of mine needed access to some sites also "protected" by this method. As I had free time in those days, I went slightly overboard in attempting to create a nice, scalable and tested framework which I could quickly iterate on.

The general idea is that we solve the "proof of work" anti-bot framework asyncronously from the crawler that is collecting data. This token seems to be valid for a long period of time, assuming the ip address does not jump around much. We do this by mocking out a JS document and letting the JS "just run" in the sandbox. We've hooked everything that it requires and force feed it "good" data.

## This doesn't work anymore! / I want support

It did, since ~2017. At the time of release this, it still worked. Though I don't utilize it, nor do any of my clients, so I don't intend to support it anyway. That is why it is open sourced.

Want hints? I'd start by running the tests and upgrading dependancies!

## License

```
Red Naga / Tim Strazzere (c) 2018 - *

GNU GENERAL PUBLIC LICENSE

Version 3, 29 June 2007

Copyright (C) 2007 Free Software Foundation, Inc.
<https://fsf.org/>

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.
```