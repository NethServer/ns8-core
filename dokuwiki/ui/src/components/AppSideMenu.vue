<template>
  <cv-side-nav id="side-nav" fixed expanded class="app-side-nav">
    <AppSideMenuContent />
  </cv-side-nav>
</template>

<script>
import AppSideMenuContent from "@/components/AppSideMenuContent";

export default {
  name: "AppSideMenu",
  components: { AppSideMenuContent },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

// @media (max-width: $breakpoint-large) { ////
//   // hide side menu on small devices
//   #side-nav {
//     display: none;
//   }
// }
</style>
