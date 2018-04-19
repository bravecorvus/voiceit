<template>
  <div>
    <Login v-if="selection == 'login'"/>
    <Register v-else-if="selection == 'register'"/>
    <div id="maindiv" v-else>
      <br />
      <br />
      <br />
      <br />
      <br />
      <br />
      <br />
      <button @click="changeselection('login')">Login</button>
      <button @click="changeselection('register')">Register</button>
    </div>
  </div>
</template>

<script>
import Login from './Login';
import Register from './Register';

export default {
  name: 'Main',
  components: {
    Login,
    Register,
  },
  data() {
    return {
      selection: null,
    };
  },
  methods: {
    changeselection(arg) {
      this.selection = arg;
    },
  },
  mounted() {
    // See if user is logged in (from browser cookies) by asking the /isloggedin endpoint
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="sass" scoped>
  #maindiv
    text-align: center
</style>
