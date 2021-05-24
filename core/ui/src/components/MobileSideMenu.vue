<template>
  <transition name="slide-menu">
    <div
      v-if="isMenuShown"
      class="mobile-side-menu cv-side-nav bx--side-nav bx--side-nav__navigation bx--side-nav--expanded"
    >
      <!-- v-click-outside="clickOutside" //// -->
      <SideMenuContent />
    </div>
  </transition>
</template>

<script>
import SideMenuContent from "@/components/SideMenuContent";

export default {
  name: "MobileSideMenu",
  components: { SideMenuContent },
  data() {
    return {
      isMenuShown: false,
      isClickOutsideEnabled: false,
    };
  },
  created() {
    // register to logout event
    this.$root.$on("toggleMobileSideMenu", this.toggleMobileSideMenu);
  },
  mounted() {
    // prevent glitch: click-outside is incorrectly detected when mobile side menu appears
    setTimeout(() => {
      this.isClickOutsideEnabled = true;
    }, 200);
  },
  beforeDestroy() {
    // remove event listener
    this.$root.$off("toggleMobileSideMenu", this.toggleMobileSideMenu);
  },
  methods: {
    toggleMobileSideMenu() {
      this.isMenuShown = !this.isMenuShown;
    },
    clickOutside() {
      console.log("click outside mobile side menu", this.isMenuShown); //// fix

      setTimeout(() => {
        this.isMenuShown = false;
      }, 200);
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.mobile-side-menu {
  // background-color: $ui-05; ////
  // border-left: 1px solid $interactive-02; ////
  // color: $ui-01; ////
  // height: 10rem; ////
  width: $side-menu-width;
  // min-width: 23rem; ////
  height: calc(100vh - 3rem);
  position: fixed;
  top: 3rem;
  left: 0;
  z-index: 10000;
  overflow: auto;
  // padding: 1rem; ////
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
