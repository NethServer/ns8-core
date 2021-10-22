<template>
  <div
    data-overflow-menu
    :class="`cv-overflow-menu ${carbonPrefix}--overflow-menu`"
    :id="uid"
  >
    <button
      :class="[
        `${carbonPrefix}--overflow-menu__trigger ${carbonPrefix}--tooltip__trigger`,
        `${carbonPrefix}--tooltip--a11y`,
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
      <span :class="`${carbonPrefix}--assistive-text`" v-if="label">{{
        label
      }}</span>

      <slot name="trigger">
        <OverflowMenuVertical16
          :class="`${carbonPrefix}--overflow-menu__icon`"
        />
      </slot>
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
import OverflowMenuVertical16 from "@carbon/icons-vue/es/overflow-menu--vertical/16";
// import { uidMixin, carbonPrefixMixin, methodsMixin } from "../../mixins"; ////

export default {
  name: "HeaderGlobalMenu",
  extends: CvOverflowMenu,
  components: { OverflowMenuVertical16 },
  //   mixins: [
  //     uidMixin,
  //     carbonPrefixMixin,
  //     methodsMixin({ trigger: ["blur", "focus"] }),
  //   ],
  props: {
    label: String,
    flipMenu: Boolean,
    up: Boolean,
    offset: {
      type: Object,
      validator(value) {
        return value && value.left !== undefined && value.top !== undefined;
        // value.hasOwnProperty('left') && value.hasOwnProperty('top');
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
