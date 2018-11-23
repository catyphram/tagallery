FROM nginx:alpine

# Install node and yarn
RUN apk add --update make gnupg curl nodejs-current yarn

# Copy and build source files
COPY . /tmp/client
WORKDIR /tmp/client/
RUN yarn install && yarn build

RUN cp -r dist/. /usr/share/nginx/html/