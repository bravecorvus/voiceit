<template>
  <div id="register">
    <link href="http://vjs.zencdn.net/6.6.3/video-js.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/videojs-record/2.1.3/css/videojs.record.min.css" rel="stylesheet">

    <h1>Register</h1>
    <p v-if="showintro">Will record 3 sets of videos
    (please say "My face and my voice identify me.").
    Start talking when you the video pops up on the screen.
    When you are ready, please enter the username you want to use and press submit.</p>
    <p v-if="nextvideo">Recording next video in {{ this.nextvideocounter }}</p>
    <div style="align: center;" id="recorddiv">
      <p v-if="showcountdown">Recording video for {{countdown}} seconds.</p>
      <video style="align: center !important;"
        id="auth-video" class="video-js vjs-default-skin"></video>
    </div>

    <div id="username">
      <br />
      <br />
      <br />
      <br />
      <label>Username</label>
      <input v-model="username" type="text"/>
      <button @click="submit">Submit</button>
    </div>

    <div id="processing">
      <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css">
      <br />
      <br />
      <br />
      <br />
      <p>Attempting to Register</p>
      <i style="font-size: 50px; color: black;" class="fa fa-spinner fa-spin" id="spinner"></i>
    </div>
    <div id="err">
    </div>
  </div>
</template>

<script>
require('recordrtc');
require('webrtc-adapter/out/adapter');
const videojs = require('video.js');
require('videojs-record');
const is = require('is_js');
const axios = require('axios');


export default {
  name: 'Register',
  data() {
    return {
      username: '',
      countdown: 5,
      showcountdown: false,
      player: null,
      showintro: true,
      counter: 2,
      fordata: null,
      nextvideo: false,
      nextvideocounter: 3,
    };
  },
  methods: {
    submit() {
      $('#username').css('display', 'none');

      function sleep(time) {
        return new Promise(resolve => setTimeout(resolve, time));
      }

      $('#username').css('display', 'none');
      $('#processing').css('display', 'none');
      $('#recorddiv').css('display', 'none');


      this.player = videojs('auth-video', {
      // video.js options
        controls: true,
        plugins: {
        // videojs-record plugin options
          record: {
            audio: true,
            video: true,
            maxLength: 5,
            videoMimeType: 'video/mp4;codecs=H264',
          },
        },
      });

      this.player.on('deviceError', () => {
        console.log('device error:', this.player.deviceErrorCode);
      });

      this.player.on('error', (error) => {
        console.log('error:', error);
      });

      this.player.on('finishRecord', () => {
        if (this.counter < 4) {
          $('#recorddiv').css('display', 'none');
          if (is.firefox()) {
            this.formdata.append(`file${this.counter}`, this.player.recordedData, this.username + this.counter);
          } else if (is.chrome()) {
            this.formdata.append(`file${this.counter}`, this.player.recordedData.video, this.username + this.counter);
          } else {
            this.formdata.append(`file${this.counter}`, this.player.recordedData, this.username + this.counter);
          }
          this.counter += 1;
          this.nextvideo = true;

          sleep(1000).then(() => {
            this.nextvideocounter -= 1;
            sleep(1000).then(() => {
              this.nextvideocounter -= 1;
              sleep(1000).then(() => {
                this.nextvideocounter -= 1;
                sleep(1000).then(() => {
                  this.nextvideocounter = 3;
                  this.nextvideo = false;
                  this.player.record().reset();
                  this.player.deviceButton.trigger('click');
                  this.countdown = 5;
                  $('#recorddiv').css('display', '');
                  $('.vjs-control-bar').css('display', 'none');
                  $('.vjs-record-button').trigger('click');
                  this.countdown -= 1;
                  this.showcountdown = true;
                  sleep(1000).then(() => {
                    this.countdown -= 1;
                    sleep(1000).then(() => {
                      this.countdown -= 1;
                      sleep(1000).then(() => {
                        this.countdown -= 1;
                        sleep(1000).then(() => {
                          this.recording = false;
                          this.showcountdown = true;
                        });
                      });
                    });
                  });
                });
              });
            });
          });
        } else {
          if (is.firefox()) {
            this.formdata.append('file', this.player.recordedData, this.username);
          } else if (is.chrome()) {
            this.formdata.append('file', this.player.recordedData.video, this.username);
          } else {
            this.formdata.append('file', this.player.recordedData, this.username);
          }
          $('#recorddiv').css('display', 'none');
          $('#processing').css('display', '');

          axios.post(
            '/register',
            this.formdata,
            {
              headers: { 'Content-Type': 'multipart/form-data' },
            },
          ).then(() => {
            $('#err').prepend('<p style="text-align: center;">Successfully created account. Please try logging in after we redirect you in a few seconds</p>');
            sleep(3000).then(() => {
              window.location.reload();
            });
            $('#processing').css('display', 'none');
            this.player.record().destroy();
          })
            .catch(() => {
              $('#err').prepend('<p style="text-align: center;">Error creating account. Redirecting back to home page in a few seconds.</p>');
              $('#processing').css('display', 'none');
              sleep(3000).then(() => {
                window.location.reload();
              });
            });
        }
      });

      this.permissions();
      this.formdata = new FormData();
    },

    permissions() {
      this.player.deviceButton.trigger('click');

      if (navigator.mediaDevices === undefined) {
        navigator.mediaDevices = {};
      }
      // Some browsers partially implement mediaDevices. We can't just assign an object
      // with getUserMedia as it would overwrite existing properties.
      // Here, we will just add the getUserMedia property if it's missing.
      if (navigator.mediaDevices.getUserMedia === undefined) {
        navigator.mediaDevices.getUserMedia = ((constraints) => {
          // First get ahold of the legacy getUserMedia, if present
          const getUserMedia = navigator.webkitGetUserMedia || navigator.mozGetUserMedia;
          // Some browsers just don't implement it - return a rejected promise with an error
          // to keep a consistent interface
          if (!getUserMedia) {
            return Promise.reject(new Error('getUserMedia is not implemented in this browser'));
          }

          // Otherwise, wrap the call to the old navigator.getUserMedia with a Promise
          return new Promise(((resolve, reject) => {
            getUserMedia.call(navigator, constraints, resolve, reject);
          }));
        });
      }

      navigator.mediaDevices.getUserMedia({ audio: true, video: true })
        .then(() => {
          this.permissionsuccess();
        })
        .catch((err) => {
          console.log(`${err.name}: ${err.message}`);
        });
    },

    permissionsuccess() {
      function sleep(time) {
        return new Promise(resolve => setTimeout(resolve, time));
      }


      sleep(3000).then(() => {
        this.showintro = false;
        this.showcountdown = false;
        $('.vjs-record-button').trigger('click');
        $('.vjs-control-bar').css('display', 'none');
        $('#recorddiv').css('display', '');
        sleep(1000).then(() => {
          this.showcountdown = true;
          this.countdown -= 1;
          sleep(1000).then(() => {
            this.countdown -= 1;
            sleep(1000).then(() => {
              this.countdown -= 1;
              sleep(1000).then(() => {
                this.countdown -= 1;
                sleep(1000).then(() => {
                  this.recording = false;
                  this.showcountdown = false;
                });
              });
            });
          });
        });
      });
    },
  },

  mounted() {
    $('#recorddiv').css('display', 'none');
    $('#processing').css('display', 'none');
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="sass" scoped>
#register
  text-align: center
h1
  font-size: 30px
  text-align: center
</style>
