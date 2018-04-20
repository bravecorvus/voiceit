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

FROM golang:1.10-alpine as backend-build-env
WORKDIR /go/src/github.com/gilgameshskytrooper/voiceit/
RUN apk --no-cache add ca-certificates && apk --no-cache add git
COPY /backend ./backend
RUN go get -d -v ./...
RUN cd backend && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o voiceit .

FROM jrottenberg/ffmpeg:3.3-scratch as ffmpeg-build-env

FROM scratch
COPY --from=ffmpeg-build-env / .
COPY --from=frontend-build-env /usr/app/dist /dist
COPY --from=backend-build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=backend-build-env /go/src/github.com/gilgameshskytrooper/voiceit/backend/voiceit /
CMD ["./voiceit"]
