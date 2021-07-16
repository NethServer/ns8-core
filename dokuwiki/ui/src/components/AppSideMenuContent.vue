<template>
  <div>
    <div class="module-id">
      <span>app name ////</span>
    </div>

    <cv-side-nav-items>
      <cv-side-nav-link
        @click="goToPage('status')"
        :class="{ 'current-page': isLinkActive('status') }"
      >
        <template v-slot:nav-icon><Activity20 /></template>
        <span>{{ $t("status.title") }}</span>
      </cv-side-nav-link>
      <cv-side-nav-link
        @click="goToPage('settings')"
        :class="{ 'current-page': isLinkActive('settings') }"
      >
        <template v-slot:nav-icon><Settings20 /></template>
        <span>{{ $t("settings.title") }}</span>
      </cv-side-nav-link>
      <cv-side-nav-link
        @click="goToPage('logs')"
        :class="{ 'current-page': isLinkActive('logs') }"
      >
        <template v-slot:nav-icon><Catalog20 /></template>
        <span>{{ $t("logs.title") }}</span>
      </cv-side-nav-link>
      <cv-side-nav-link
        @click="goToPage('about')"
        :class="{ 'current-page': isLinkActive('about') }"
      >
        <template v-slot:nav-icon><Information20 /></template>
        <span>{{ $t("about.title") }}</span>
      </cv-side-nav-link>
    </cv-side-nav-items>
  </div>
</template>

<script>
import Settings20 from "@carbon/icons-vue/es/settings/20";
import Catalog20 from "@carbon/icons-vue/es/catalog/20";
import Information20 from "@carbon/icons-vue/es/information/20";
import Activity20 from "@carbon/icons-vue/es/activity/20";

let nethserver = window.nethserver;

export default {
  name: "AppSideMenuContent",
  components: {
    Catalog20,
    Settings20,
    Information20,
    Activity20,
  },
  created() {
    // register to appNavigation event
    this.$root.$on("appNavigation", this.onAppNavigation);
  },
  beforeDestroy() {
    // remove event listener
    this.$root.$off("appNavigation");
  },
  methods: {
    isLinkActive(page) {
      return nethserver.getPage() === page;
    },
    //// move to mixin
    goToPath(path) {
      if (window.parent.ns8.$route.fullPath != path) {
        window.parent.ns8.$router.push(path);
      }
      this.$forceUpdate();
    },
    goToPage(page) {
      //// parametrize dokuwiki1
      const path = "/apps/dokuwiki1?page=" + page;

      console.log("goToPage", path, window.parent.ns8.$route.fullPath); ////

      if (window.parent.ns8.$route.fullPath != path) {
        window.parent.ns8.$router.push(path);
      }
    },
    onAppNavigation(page) {
      console.log("onAppNavigation", page); ////
      this.$forceUpdate();
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.module-id {
  height: 3rem;
  padding-left: 1.5rem;
  margin-top: 1rem;
  display: flex;
  align-items: center;
  font-weight: 600;
}
</style>
