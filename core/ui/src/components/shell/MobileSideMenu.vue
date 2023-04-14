<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <transition name="slide-menu">
    <div
      v-if="isMobileSideMenuShown"
      v-click-outside="clickOutside"
      class="core-side-nav mobile-side-menu cv-side-nav bx--side-nav bx--side-nav__navigation bx--side-nav--expanded"
    >
      <SideMenuContent />
    </div>
  </transition>
</template>

<script>
import SideMenuContent from "@/components/shell/SideMenuContent";
import { mapState, mapActions } from "vuex";

export default {
  name: "MobileSideMenu",
  components: { SideMenuContent },
  data() {
    return {
      isTransitioning: false,
    };
  },
  computed: {
    ...mapState(["isMobileSideMenuShown"]),
  },
  watch: {
    isMobileSideMenuShown: function () {
      this.isTransitioning = true;

      setTimeout(() => {
        this.isTransitioning = false;
      }, 300); // same duration as .slide-menu transition
    },
  },
  methods: {
    ...mapActions(["setMobileSideMenuShownInStore"]),
    clickOutside() {
      if (!this.isTransitioning) {
        // close menu
        this.setMobileSideMenuShownInStore(false);
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.mobile-side-menu {
  width: $side-menu-width;
  height: calc(100vh - 3rem);
  position: fixed;
  top: 3rem;
  left: 0;
  z-index: 10000;
  overflow: auto;
}

.slide-menu-enter-active,
.slide-menu-leave-active {
  transition: all 0.3s ease;
}

.slide-menu-enter,
.slide-menu-leave-to {
  transform: translateX(-$side-menu-width);
}
</style>
