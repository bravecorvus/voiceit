# VoiceIt Web App
> Andrew Lee

## Introduction

This web application allows for user registration (using video enrollment) as well as user authentication (via video verification).

I am utilizing [Vue.js](http://vuejs.org/) as a basis for the frontend (as well as heavy use of the [videojs-record](https://github.com/collab-project/videojs-record) project to capture video and audio from users in the client side and sending the resulting blob to the server). The backend is written in [Go](https://golang.org/), and I am using [Redis](https://redis.io/) to store persistent data.

## Build Instructions
In order to minimize the system dependencies, I have used [docker-compose](https://docs.docker.com/compose/) to build the frontend and backend, and run that plus the Redis database all in one command. Because of this, the only system dependencies to run this program is `docker version > 18.0.3.0-ce` and `docker-compose version > 1.20.1`.

I am using a new technique called "statically compiled builds", which defines multiple build steps in a single Dockerfile to compile (in the case of my Go backend server), or produce pure client side code (in the case of the Vue frontend), and copy those assets into a [Docker Scratch](https://hub.docker.com/_/scratch/) image (which is described by Docker as an explicitly empty image). This is possible because the resulting code has no system dependencies, and can be run directly in an empty environment (a testament to the power of Go). Furthermore, the resulting image would have been about `14MB`. However, due to the limitations of native browser support for encoding video (Chrome only exporting Matroska H264 and Firefox exporting WebM videos), I also had to include a statically compiled executable of the [ffmpeg](https://hub.docker.com/r/jrottenberg/ffmpeg/) program in order process the video and encode it into full H264/MPEG-4 AVC format on the web server (which is required by the VoiceIt API. This increases the total size to be `63.5MB`.

![dockerimages](https://78.media.tumblr.com/62023c3577c61cd20fa1e82b3be2ebf6/tumblr_p7h5ne4MwC1s5a4bko1_1280.png)

The only thing that needs to be done in order to run the program is to first copy the file `docker-compose.yaml` to be `docker-compose-actual.yaml`. This is done so that you can store the actual secret in the environment variables specified in `docker-compose-actual.yaml` without worrying about exposing secrets such as VoiceIt API key and value since it is added to `.gitignore`.

You will need to replace the variables `VOICEITAPIKEY` and `VOICEITAPITOKEN` in `docker-compose-actual.yaml` with real values provided by VoiceIt.
```
    environment:
      - VOICEITAPIKEY=key_111111111111111111111111
      - VOICEITAPITOKEN=tok_11111111111111111111111
```


Then, just run the following command:
```
./run-dev
```

**Note: I have noticed that running the above command in a system which defaults to the root user will not work since the Redis container needs to be owned by user:group `1001:1001` (due to their implementation of Linxu security). Hence, you will need to add the line `chown -R 1001:1001 redis` after the `mkdir` line in the above script for this scenario**

This programs bootstraps the directories we will be mounting on the host system as well as moving the files `docker-compose.yaml` and `docker-compose-actual.yaml` so that `actual` gets used during execution, but switches back places after the program is shut down via `CTRL-C`.

**Please make sure you run this program AT LEAST ONCE even if you already understand how to use `docker-compose` since the directories are necessary to run Redis and the web server.**

![run-dev](https://78.media.tumblr.com/6c7037e4dd74023d34213aac24002ac3/tumblr_p7h7ljqUCv1s5a4bko1_1280.png)

Once you run the above program at least once, and you are sure you will no longer be updating the git remote repo (you no longer need to worry about uploading your API key and value to Github), you can replace the original `docker-compose.yaml` with `docker-compose-actual.yaml`, and run the container headlessly via the command `docker-compose up -d`.

Furthermore, this will allow you to define scripts to ensure `docker-compose up` is run after every restart.

Heres an example of how you could accomplish that in `systemd`:
File: `/etc/systemd/system/voiceit.service`
```
[Unit]
Description=VoiceIt docker-compose container starter
After=docker.service network-online.target
Requires=docker.service network-online.target

[Service]
WorkingDirectory=/root/go/src/github.com/gilgameshskytrooper/voiceit
Type=oneshot
RemainAfterExit=yes

ExecStart=/usr/bin/docker-compose up -d

ExecStop=/usr/local/bin/docker-compose down

ExecReload=/usr/bin/docker-compose up -d

[Install]
WantedBy=multi-user.target
```
**Just change `WorkingDirectory` accordingly**

## Caveats
I was having some amount of trouble getting video verification to work, so I took some liberties to fully showcase the application functionality even if video verification was not working for me at this time.

Namely, registration will work correctly, but I have not been able to login using video verification. Hence, I have a counter in the backend which automatically logs you in (even with a failed API video verification response from VoiceIt) every other try.

Furthermore, in order to showcase the server side video processing capacity, I am not explicitly deleting any uploaded or converted videos. In a true production setting, both these feature showcasing functionality will be turned off since a authentication which doesn't authenticate 50% of the time is useless, and videos after creating a new video enrollment or authenticating a user can be deleted to save space as they are only needed for the duration of the VoiceIt API call.
