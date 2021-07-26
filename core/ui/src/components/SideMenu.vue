<template>
  <cv-side-nav
    id="side-nav"
    :fixed="!isExternalApp"
    :expanded="!isExternalApp"
    :rail="isExternalApp"
  >
    <SideMenuContent />
  </cv-side-nav>
</template>

<script>
import SideMenuContent from "@/components/SideMenuContent";

export default {
  name: "SideMenu",
  components: { SideMenuContent },
  computed: {
    isExternalApp() {
      return this.$route.path.includes("/apps/");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

@media (max-width: $breakpoint-large) {
  // hide side menu on small devices
  #side-nav {
    display: none;
  }
}
</style>
