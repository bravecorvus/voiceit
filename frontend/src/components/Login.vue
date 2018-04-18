<template>
  <div>
    <link href="http://vjs.zencdn.net/6.6.3/video-js.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/videojs-record/2.1.3/css/videojs.record.min.css" rel="stylesheet">
    <p>Login</p>
    <video id="auth-video" class="video-js vjs-default-skin"></video>
  </div>
</template>

<script>
require('recordrtc');
require('webrtc-adapter/out/adapter');
const videojs = require('video.js');
require('videojs-record');


export default {
  name: 'Login',
  data() {
    return {
      segmentNumber: 0,
    };
  },
  mounted() {
    const player = videojs('auth-video', {
      // video.js options
      // controls: true,
      width: 1280,
      height: 720,
      fluid: true,
      plugins: {
        // videojs-record plugin options
        record: {
          audio: true,
          video: true,
          // maxLength: 10,
          frameWidth: 1280,
          frameHeight: 720,
        },
      },
    });

    player.on('deviceError', () => {
      console.log('device error:', player.deviceErrorCode);
    });

    player.on('error', (error) => {
      console.log('error:', error);
    });

    // user clicked the record button and started recording
    player.on('startRecord', () => {
      console.log('started recording!');
    });

    // player.on('timestamp', () => {
    //       if (player.recordedData && player.recordedData.length > 0) {
    //         const binaryData = player.recordedData[player.recordedData.length - 1];
    //
    //         this.segmentNumber += 1;

    //         const formData = new FormData();
    //         formData.append('SegmentNumber', this.segmentNumber);
    //         formData.append('Data', binaryData);

    //         $.ajax({
    //           url: '/login',
    //           method: 'POST',
    //           data: formData,
    //           cache: false,
    //           processData: false,
    //           contentType: false,
    //           success() {
    //             console.log(`segment: ${this.segmentNumber}`);
    //           },
    //         });
    //       }
    //     });

    // user completed recording and stream is available
    player.on('finishRecord', () => {
    // the blob object contains the recorded data that
    // can be downloaded by the user, stored on server etc.
      console.log('finished recording: ', player.recordedData);
    });

    // Start recording immedately after loading the page
    $('.vjs-device-button').trigger('click');
    $.ajax({
      url: '/login',
      method: 'POST',
    });
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="sass" scoped>
</style>
