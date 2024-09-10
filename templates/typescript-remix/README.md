# Welcome to Remix!

[![Deploy on Jetify](https://www.jetify.com/img/devbox/deploy-on-jetify-galaxy.svg)](https://cloud.jetify.com/team/new/deploy?repo=github.com/jetify-examples/typescript-remix)

- 📖 [Remix docs](https://remix.run/docs)

## Development

Run the dev server:

```shellscript
npm run dev
```

## Deployment

First, build your app for production:

```sh
npm run build
```

Then run the app in production mode:

```sh
npm start
```

Now you'll need to pick a host to deploy it to.

### DIY

If you're familiar with deploying Node applications, the built-in Remix app server is production-ready.

Make sure to deploy the output of `npm run build`

- `build/server`
- `build/client`
