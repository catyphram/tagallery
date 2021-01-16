FROM node:alpine AS builder

RUN apk add --update --no-cache python make g++

WORKDIR /app
COPY package.json yarn.lock ./

RUN yarn install --frozen-lockfile --non-interactive
RUN yarn cache clean

COPY . ./

RUN yarn build
RUN yarn --frozen-lockfile --non-interactive --production

FROM node:alpine

WORKDIR /app

ENV NODE_ENV=production

COPY package.json nuxt.config.js tsconfig.json ./
COPY --from=builder ./app/node_modules ./node_modules/
COPY --from=builder ./app/.nuxt ./.nuxt/
COPY --from=builder ./app/static ./static/

CMD [ "yarn", "start" ]
