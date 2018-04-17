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
COPY /backend/main.go .
COPY /backend/email ./email
COPY /backend/structs ./structs
COPY /backend/utils ./utils
COPY /backend/app ./app
RUN go get -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o voiceit .

FROM scratch
COPY --from=frontend-build-env /usr/app/dist /dist
COPY --from=backend-build-env /tmp /tmp
COPY --from=backend-build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=backend-build-env /go/src/github.com/gilgameshskytrooper/voiceit/voiceit /
CMD ["./voiceit"]
