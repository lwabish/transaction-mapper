<template>
  <div id="app">
    <h1>Transaction Mapper Web</h1>
    <FileUploadForm />
  </div>
</template>

<script setup>
import FileUploadForm from './components/form.vue'
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>