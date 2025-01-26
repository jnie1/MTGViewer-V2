<script setup lang="ts">
import {RouterLink} from 'vue-router';
</script>

<template>
    <nav>
      <ul>
        <li><RouterLink to="/">Home</RouterLink></li>
        <li><RouterLink to="/about">About</RouterLink></li>
        <li><RouterLink to="/card">Card</RouterLink></li>
        <li><RouterLink to="/signup">SignUp</RouterLink></li>
      </ul> 
    </nav>
</template>


<style scoped>
  nav {
    color:antiquewhite;
    position: relative;
  }

  ul, li {
    list-style-type: none;
    display: inline;
    padding: 1%;
    font-size: 15px;
  }
</style>
