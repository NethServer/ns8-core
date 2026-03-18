<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <cv-side-nav
    id="side-nav"
    :fixed="!isExternalApp"
    :expanded="!isExternalApp"
    :rail="isExternalApp"
    class="core-side-nav"
  >
    <SideMenuContent />
  </cv-side-nav>
</template>

<script>
import SideMenuContent from "@/components/shell/SideMenuContent";

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
@import "../../styles/carbon-utils";

@media (max-width: $breakpoint-large) {
  // hide side menu on small devices
  #side-nav {
    display: none;
  }
}
</style>

<style lang="scss">
@import "../../styles/carbon-utils";

// global styles

.core-side-nav.bx--side-nav {
  background-color: $ui-05;
  color: $ui-01;
}

.core-side-nav .bx--side-nav__submenu {
  color: $ui-01;
}

.core-side-nav a.bx--side-nav__link > .bx--side-nav__link-text,
.core-side-nav.bx--side-nav a.bx--header__menu-item .bx--text-truncate-end {
  color: $ui-01;
}

.core-side-nav .bx--side-nav__menu a.bx--side-nav__link--current > span,
.core-side-nav
  .bx--side-nav__menu
  a.bx--side-nav__link[aria-current="page"]
  > span,
.core-side-nav a.bx--side-nav__link--current > span {
  color: $ui-01;
}

.core-side-nav .bx--side-nav__icon > svg {
  width: 20px;
  height: 20px;
  fill: $ui-01;
}

.core-side-nav .bx--side-nav__menu a.bx--side-nav__link--current,
.core-side-nav .bx--side-nav__menu a.bx--side-nav__link[aria-current="page"],
.core-side-nav a.bx--side-nav__link--current {
  background-color: $interactive-02;
}

.core-side-nav .bx--side-nav__submenu:hover {
  color: $ui-01;
  background-color: $interactive-02;
}

.core-side-nav
  .bx--side-nav__item:not(.bx--side-nav__item--active):hover
  .bx--side-nav__item:not(.bx--side-nav__item--active)
  > .bx--side-nav__submenu:hover,
.core-side-nav
  .bx--side-nav__item:not(.bx--side-nav__item--active)
  > .bx--side-nav__link:hover,
.core-side-nav
  .bx--side-nav__menu
  a.bx--side-nav__link:not(.bx--side-nav__link--current):not([aria-current="page"]):hover,
.core-side-nav .bx--side-nav a.bx--header__menu-item:hover,
.core-side-nav
  .bx--side-nav
  .bx--header__menu-title[aria-expanded="true"]:hover {
  color: $ui-01;
  background-color: $interactive-02;
}

.core-side-nav
  .bx--side-nav__item:not(.bx--side-nav__item--active)
  > .bx--side-nav__link:hover
  > span,
.core-side-nav
  .bx--side-nav__item:not(.bx--side-nav__item--active)
  .bx--side-nav__menu-item
  > .bx--side-nav__link:hover
  > span {
  color: $ui-01;
}

.core-side-nav a.bx--side-nav__link[aria-current="page"]::before,
.core-side-nav a.bx--side-nav__link--current::before {
  background-color: $inverse-support-04;
}

.core-side-nav .bx--side-nav__submenu {
  height: 3rem;
}

.core-side-nav .cv-side-nav-item-link.bx--side-nav__link {
  height: 3rem;
}

.core-side-nav a.bx--side-nav__link[aria-current="page"],
.core-side-nav a.bx--side-nav__link--current {
  background-color: $inverse-02;
}

.core-side-nav
  a.bx--side-nav__link[aria-current="page"]
  .bx--side-nav__link-text,
.core-side-nav a.bx--side-nav__link--current .bx--side-nav__link-text {
  color: $ui-01;
}

.core-side-nav .bx--side-nav--rail {
  width: $side-menu-collapsed-width;
}

// menu divider
.core-side-nav .bx--side-nav__divider {
  background-color: $inverse-02;
}
</style>
