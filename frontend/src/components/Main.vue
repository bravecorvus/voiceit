<template>
  <div>
    <Home v-if="loggedin"/>
    <Login v-else/>
  </div>
</template>

<script>
import Home from './Home';
import Login from './Login';

export default {
  name: 'Main',
  data() {
    return {
      loggedin: false,
    };
  },
  components: {
    Home,
    Login,
  },
  mounted() {
    // See if user is logged in (from browser cookies) by asking the /isloggedin endpoint
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="sass" scoped>
</style>
