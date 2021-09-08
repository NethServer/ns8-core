<template>
  <cv-side-nav
    id="side-nav"
    :fixed="!isExternalApp"
    :expanded="!isExternalApp"
    :rail="isExternalApp"
    class="ns8-side-nav"
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

<style lang="scss">
@import "../styles/carbon-utils";

// global styles

.ns8-side-nav.bx--side-nav {
  background-color: $ui-05;
  color: $ui-01;
}

.ns8-side-nav .bx--side-nav__submenu {
  color: $ui-01;
}

.ns8-side-nav a.bx--side-nav__link > .bx--side-nav__link-text,
.ns8-side-nav.bx--side-nav a.bx--header__menu-item .bx--text-truncate-end {
  color: $ui-01;
}

.ns8-side-nav .bx--side-nav__menu a.bx--side-nav__link--current > span,
.ns8-side-nav
  .bx--side-nav__menu
  a.bx--side-nav__link[aria-current="page"]
  > span,
.ns8-side-nav a.bx--side-nav__link--current > span {
  color: $ui-01;
}

.ns8-side-nav .bx--side-nav__icon > svg {
  width: 20px;
  height: 20px;
  fill: $ui-01;
}

.ns8-side-nav .bx--side-nav__menu a.bx--side-nav__link--current,
.ns8-side-nav .bx--side-nav__menu a.bx--side-nav__link[aria-current="page"],
.ns8-side-nav a.bx--side-nav__link--current {
  background-color: $interactive-02;
}

.ns8-side-nav .bx--side-nav__submenu:hover {
  color: $ui-01;
  background-color: $interactive-02;
}

.ns8-side-nav
  .bx--side-nav__item:not(.bx--side-nav__item--active):hover
  .bx--side-nav__item:not(.bx--side-nav__item--active)
  > .bx--side-nav__submenu:hover,
.ns8-side-nav
  .bx--side-nav__item:not(.bx--side-nav__item--active)
  > .bx--side-nav__link:hover,
.ns8-side-nav
  .bx--side-nav__menu
  a.bx--side-nav__link:not(.bx--side-nav__link--current):not([aria-current="page"]):hover,
.ns8-side-nav .bx--side-nav a.bx--header__menu-item:hover,
.ns8-side-nav
  .bx--side-nav
  .bx--header__menu-title[aria-expanded="true"]:hover {
  color: $ui-01;
  background-color: $interactive-02;
}

.ns8-side-nav
  .bx--side-nav__item:not(.bx--side-nav__item--active)
  > .bx--side-nav__link:hover
  > span,
.ns8-side-nav
  .bx--side-nav__item:not(.bx--side-nav__item--active)
  .bx--side-nav__menu-item
  > .bx--side-nav__link:hover
  > span {
  color: $ui-01;
}

.ns8-side-nav a.bx--side-nav__link[aria-current="page"]::before,
.ns8-side-nav a.bx--side-nav__link--current::before {
  background-color: $inverse-support-04;
}

.ns8-side-nav .bx--side-nav__submenu {
  height: 3rem;
}

.ns8-side-nav .cv-side-nav-item-link.bx--side-nav__link {
  height: 3rem;
}

.ns8-side-nav a.bx--side-nav__link[aria-current="page"],
.ns8-side-nav a.bx--side-nav__link--current {
  background-color: $inverse-02;
}

.ns8-side-nav
  a.bx--side-nav__link[aria-current="page"]
  .bx--side-nav__link-text,
.ns8-side-nav a.bx--side-nav__link--current .bx--side-nav__link-text {
  color: $ui-01;
}

.ns8-side-nav .bx--side-nav--rail {
  width: $side-menu-collapsed-width;
}
</style>
