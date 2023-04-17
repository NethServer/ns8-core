<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div
    data-overflow-menu
    :class="[`cv-overflow-menu`, 'header-global-menu']"
    :id="uid"
  >
    <button
      :class="[
        `${carbonPrefix}--overflow-menu__trigger ${carbonPrefix}--tooltip__trigger`,
        `${carbonPrefix}--tooltip--a11y`,
        'bx--btn--primary',
        'bx--btn',
        'bx--btn--icon-only',
        'bx--header__action',
        {
          [`${this.carbonPrefix}--tooltip--${tipPosition}`]: label,
          [`${this.carbonPrefix}--tooltip--align-${tipAlignment}`]: label,
          [`${carbonPrefix}--overflow-menu--open`]: open,
        },
      ]"
      aria-haspopup
      type="button"
      :aria-expanded="open ? 'true' : 'false'"
      :aria-controls="`${uid}-menu`"
      :id="`${uid}-trigger`"
      ref="trigger"
      @click="doToggle"
      @keydown.space.prevent
      @keyup.space.prevent="doToggle"
      @keydown.enter.prevent="doToggle"
      @keydown.tab="onOverflowMenuTab"
    >
      <span :class="`${carbonPrefix}--assistive-text`" v-if="label && !open">{{
        label
      }}</span>

      <slot name="trigger"> </slot>
    </button>
    <div
      :class="[
        `${carbonPrefix}--overflow-menu-options`,
        {
          [`${carbonPrefix}--overflow-menu--flip`]: flipMenu,
          [`${carbonPrefix}--overflow-menu-options--open`]: open,
        },
      ]"
      tabindex="-1"
      ref="popup"
      :aria-labelledby="`${uid}-trigger`"
      :id="`${uid}-menu`"
      :style="{ left: left + 'px', top: top + 'px' }"
      @focusout="checkFocusOut"
      @mousedown.prevent="preventFocusOut"
    >
      <div
        class="cv-overflow-menu__before-content"
        ref="beforeContent"
        tabindex="0"
        style="position: absolute; height: 1px; width: 1px; left: -9999px"
        @focus="focusBeforeContent"
      />
      <ul :class="`${carbonPrefix}--overflow-menu-options__content`">
        <slot></slot>
      </ul>
      <div
        class="cv-overflow-menu__after-content"
        ref="afterContent"
        tabindex="0"
        style="position: absolute; height: 1px; width: 1px; left: -9999px"
        @focus="focusAfterContent"
      />
    </div>
  </div>
</template>

<script>
import { CvOverflowMenu } from "@carbon/vue";

export default {
  name: "HeaderGlobalMenu",
  extends: CvOverflowMenu,
  props: {
    label: String,
    flipMenu: Boolean,
    up: Boolean,
    offset: {
      type: Object,
      validator(value) {
        return value && value.left !== undefined && value.top !== undefined;
      },
    },
    tipPosition: {
      type: String,
      default: "right",
      validator: (val) => ["top", "left", "bottom", "right".includes(val)],
    },
    tipAlignment: {
      type: String,
      default: "center",
      validator: (val) => ["start", "center", "end"].includes(val),
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>

<style lang="scss">
@import "../../styles/carbon-utils";

// global styles

.header-global-menu .bx--tooltip__trigger svg {
  fill: currentColor !important;
}

.header-global-menu .bx--overflow-menu__trigger {
  width: 3rem;
  height: 3rem;
}
</style>
