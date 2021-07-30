<template>
  <div>
    <div class="instance-name">
      <span v-if="instanceName">{{ instanceName }}</span>
      <cv-skeleton-text
        v-else
        :width="instanceNameSkeletonWidth"
      ></cv-skeleton-text>
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
import { mapState } from "vuex";
import { QueryParamService } from "@nethserver/ns8-ui-lib";

export default {
  name: "AppSideMenuContent",
  components: {
    Catalog20,
    Settings20,
    Information20,
    Activity20,
  },
  mixins: [QueryParamService],
  data() {
    return {
      instanceNameSkeletonWidth: "70%",
    };
  },
  computed: {
    ...mapState(["instanceName", "ns8Core"]),
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
      return this.getPage() === page;
    },
    goToPath(path) {
      if (this.ns8Core.$route.fullPath != path) {
        this.ns8Core.$router.push(path);
      }
      this.$forceUpdate();
    },
    goToPage(page) {
      const path = `/apps/${this.instanceName}?page=${page}`;

      if (this.ns8Core.$route.fullPath != path) {
        this.ns8Core.$router.push(path);
      }
    },
    onAppNavigation() {
      // highlight current page in side menu
      this.$forceUpdate();
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.instance-name {
  height: 3rem;
  padding-left: 1.5rem;
  margin-top: 1rem;
  display: flex;
  align-items: center;
  font-weight: 600;
}
</style>
