<template>
  <div id="login">
    <link href="../assets/video-js.min.css" rel="stylesheet">
    <link href="../assets/videojs.record.min.css" rel="stylesheet">

    <h1 id="h1">Login</h1>
    <div style="align: center;" id="recorddiv">
      <p v-if="showvideocountdown">Recording video for {{countdown}} seconds.</p>
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
      <p>Attempting to Login</p>
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
  name: 'Login',
  data() {
    return {
      username: '',
      countdown: 5,
      player: null,
      showvideocountdown: false,
    };
  },
  methods: {
    submit() {
      $('#processing').css('display', '');
      $('#username').css('display', 'none');
      const formData = new FormData();
      if (is.firefox()) {
        formData.append('file', this.player.recordedData, this.username);
      } else if (is.chrome()) {
        formData.append('file', this.player.recordedData.video, this.username);
      } else {
        formData.append('file', this.player.recordedData, this.username);
      }

      function sleep(time) {
        return new Promise(resolve => setTimeout(resolve, time));
      }

      axios.post(
        '/login',
        formData,
        {
          headers: { 'Content-Type': 'multipart/form-data' },
        },
      ).then(() => {
        $('#h1').css('display', 'none');
        $('#processing').css('display', 'none');
        $('#err').prepend('Login successful. redirecting to secrets page in a few seconds');
        sleep(3000).then(() => {
          window.location.replace(`/secret/${this.username}`);
        });
      })
        .catch(() => {
          $('#err').prepend('<p style="text-align: center;">Failed to log in. Please try again after we redirect you back to home page in a few seconds.</p>');
          $('#processing').css('display', 'none');
          sleep(3000).then(() => {
            window.location.reload();
          });
        });
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

      sleep(100).then(() => {
        $('.vjs-record-button').trigger('click');
        $('.vjs-control-bar').css('display', 'none');
        $('#recorddiv').css('display', '');
      });
      this.showvideocountdown = true;
      sleep(1000).then(() => {
        this.countdown -= 1;
        sleep(1000).then(() => {
          this.countdown -= 1;
          sleep(1000).then(() => {
            this.countdown -= 1;
            sleep(1000).then(() => {
              this.countdown -= 1;
              sleep(1000).then(() => {
                this.recording = false;
              });
            });
          });
        });
      });
    },
  },

  mounted() {
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
      $('#username').css('display', '');
      $('#recorddiv').css('display', 'none');
    });


    this.permissions();
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="sass" scoped>
  #login
    text-align: center
  h1
    font-size: 30px
    text-align: center
</style>
