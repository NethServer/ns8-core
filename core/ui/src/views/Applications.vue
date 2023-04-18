<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
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

<style scoped lang="scss">
@import "../styles/carbon-utils";

#app-frame {
  width: calc(100% - #{$side-menu-collapsed-width});
  height: calc(100% - 3rem);
}

@media (max-width: $breakpoint-large) {
  // remove left margin since core collapsed side menu is not shown
  #app-frame {
    width: 100%;
  }
}

.iframe-embedded {
  position: absolute;
  border: none;
  z-index: 1;
}
</style>
