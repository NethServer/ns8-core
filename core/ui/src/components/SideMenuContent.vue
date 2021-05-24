<template>
  <div>
    <cv-side-nav-items>
      <cv-side-nav-menu title="Network">
        <template v-slot:nav-icon><Network_220 /></template>
        <cv-side-nav-menu-item href="javascript:void(0)" active>
          Menu item
        </cv-side-nav-menu-item>
        <cv-side-nav-menu-item href="javascript:void(0)">
          Menu item
        </cv-side-nav-menu-item>
        <cv-side-nav-menu-item href="javascript:void(0)">
          Menu item
        </cv-side-nav-menu-item>
      </cv-side-nav-menu>
      <cv-side-nav-menu title="Account provider">
        <template v-slot:nav-icon><Events20 /></template>
        <cv-side-nav-menu-item href="javascript:void(0)">
          Menu item
        </cv-side-nav-menu-item>
        <cv-side-nav-menu-item href="javascript:void(0)" aria-current="page">
          Menu item
        </cv-side-nav-menu-item>
        <cv-side-nav-menu-item href="javascript:void(0)">
          Menu item
        </cv-side-nav-menu-item>
      </cv-side-nav-menu>
      <cv-side-nav-link href="javascript:void(0)">
        <template v-slot:nav-icon><Save20 /></template>
        Backup
      </cv-side-nav-link>
      <cv-side-nav-link href="javascript:void(0)">
        <template v-slot:nav-icon><Catalog20 /></template>
        Logs
      </cv-side-nav-link>
      <cv-side-nav-link href="javascript:void(0)">
        <template v-slot:nav-icon><Settings20 /></template>
        Settings
      </cv-side-nav-link>
    </cv-side-nav-items>
  </div>
</template>

<script>
import Network_220 from "@carbon/icons-vue/es/network--2/20";
import Events20 from "@carbon/icons-vue/es/events/20";
import Save20 from "@carbon/icons-vue/es/save/20";
import Settings20 from "@carbon/icons-vue/es/settings/20";
import Catalog20 from "@carbon/icons-vue/es/catalog/20";

export default {
  name: "SideMenuContent",
  components: { Network_220, Events20, Save20, Catalog20, Settings20 },
  data() {
    return {};
  },
  methods: {},
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";
</style>
