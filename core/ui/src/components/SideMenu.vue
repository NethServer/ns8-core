<template>
  <!-- side menu -->
  <cv-side-nav id="side-nav" fixed expanded>
    <SideMenuContent />
  </cv-side-nav>
</template>

<script>
import SideMenuContent from "@/components/SideMenuContent";

export default {
  name: "SideMenu",
  components: { SideMenuContent },
  // data() {
  //   return {
  //     // mobileSideMenuShown: false, ////
  //   };
  // },
  // created() {
  //   // register to logout event
  //   this.$root.$on("toggleMobileSideMenu", this.toggleMobileSideMenu);
  // },
  // beforeDestroy() {
  //   // remove event listener
  //   this.$root.$off("toggleMobileSideMenu", this.toggleMobileSideMenu);
  // },
  // methods: {
  //   toggleMobileSideMenu() {
  //     console.log("toggleMobileSideMenu"); ////
  //     this.mobileSideMenuShown = !this.mobileSideMenuShown;
  //   },
  // },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

// .mobile-side-menu { ////
//   // background-color: $ui-05; ////
//   // border-left: 1px solid $interactive-02; ////
//   // color: $ui-01; ////
//   // height: 10rem; ////
//   width: 16rem;
//   // min-width: 23rem; ////
//   height: calc(100vh - 3rem);
//   position: fixed;
//   top: 3rem;
//   left: 0;
//   z-index: 10000;
//   overflow: auto;
//   padding: 1rem;
// }

@media (max-width: $breakpoint-large) {
  // hide side menu on small devices
  #side-nav {
    display: none;
  }
}
</style>
