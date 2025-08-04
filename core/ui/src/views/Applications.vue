<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div class="apps">
    <div>migratingApps {{ migratingApps }} ////</div>
    <div v-if="migratingApps.includes(appId)" class="migrating-app">
      <h2>{{ appId }}</h2>
      <NsInlineNotification
        kind="info"
        :title="$t('software_center.migrating_app_title')"
        :description="
          $t('software_center.migrating_app_description', {
            name: appId,
          })
        "
        :showCloseButton="false"
      />
    </div>
    <iframe
      v-else
      id="app-frame"
      class="iframe-embedded"
      :src="iframePath"
    ></iframe>
  </div>
</template>

<script>
import { mapState } from "vuex";

export default {
  name: "Applications",
  computed: {
    ...mapState(["migratingApps"]),
    appId() {
      return this.$route.params.appId;
    },
    iframePath() {
      return `apps/${this.appId}/index.html`;
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

.migrating-app {
  padding: 2rem;
}
</style>
