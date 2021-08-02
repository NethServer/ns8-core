<template>
  <cv-side-nav id="app-side-nav" fixed expanded class="app-side-nav">
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
