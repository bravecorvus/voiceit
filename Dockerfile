FROM node:9.4.0 as frontend-build-env
WORKDIR /usr/app
COPY /frontend/package.json .
COPY /frontend/yarn.lock .
COPY /frontend/index.html .
COPY /frontend/.eslintignore .
COPY /frontend/.eslintrc.js .
COPY /frontend/.postcssrc.js .
COPY /frontend/.babelrc .
COPY /frontend/build ./build
COPY /frontend/config ./config
COPY /frontend/src ./src
COPY /frontend/static ./static
RUN yarn install
RUN npm run build
RUN mkdir dist/assets
RUN wget -O /usr/app/dist/assets/video-js.min.css http://vjs.zencdn.net/6.6.3/video-js.min.css
RUN wget -O /usr/app/dist/assets/videojs.record.min.css https://cdnjs.cloudflare.com/ajax/libs/videojs-record/2.1.3/css/videojs.record.min.css

FROM golang:1.10-alpine as backend-build-env
WORKDIR /go/src/github.com/gilgameshskytrooper/voiceit/
RUN apk --no-cache add ca-certificates && apk --no-cache add git
COPY /backend ./backend
RUN go get -d -v ./...
RUN cd backend && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o voiceit .

FROM scratch
COPY --from=frontend-build-env /tmp /tmp
COPY --from=frontend-build-env /usr/app/dist /dist
COPY --from=backend-build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=backend-build-env /go/src/github.com/gilgameshskytrooper/voiceit/backend/voiceit /
COPY --from=backend-build-env /go/src/github.com/gilgameshskytrooper/voiceit/backend/templates /templates
CMD ["./voiceit"]
