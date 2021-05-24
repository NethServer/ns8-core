<template>
  <div class="apps">
    <iframe id="app-frame" class="iframe-embedded" :src="iframePath"></iframe>
  </div>
</template>

<script>
export default {
  name: "Applications",
  computed: {
    iframePath() {
      return `apps/${this.$route.params.appId}/index.html`;
    },
  },
};
</script>

<style>
#app-frame {
  width: 100%;
  height: 100%;
}
.iframe-embedded {
  position: absolute;
  width: 97%;
  height: 95%;
  border: none;
  margin-left: -20px;
  z-index: 1;
}
</style>
