<template>
  <component
    :is="tagType"
    :class="[
      `cv-tile ${carbonPrefix}--tile`,
      { [`${carbonPrefix}--tile--light`]: isLight },
      'ns-tile',
      { 'pad-bottom': icon },
    ]"
    :checked="selected"
    :expanded="expanded"
    v-bind="$attrs"
    v-on="$listeners"
    ref="tile"
  >
    <!-- can i now click on this -->
    <template>
      <slot></slot>
    </template>
    <!-- <template slot="below"> ////
      <slot name="below"></slot>
    </template> -->

    <!-- icon -->
    <NsSvg v-if="icon" :svg="icon" class="tile-icon" />
  </component>
</template>

<script>
// import CvTileClickable from "./_cv-tile-clickable"; ////
// import CvTileExpandable from "./_cv-tile-expandable";
// import CvTileSelectable from "./_cv-tile-selectable";
// import CvTileStandard from "./_cv-tile-standard";

import { CvTile } from "@carbon/vue";
import NsSvg from "./NsSvg";

export default {
  name: "NsTile",
  extends: CvTile,
  inheritAttrs: false,
  model: {
    // required for selectable kind
    prop: "modelValue",
    event: "modelEvent",
  },
  components: {
    NsSvg,
    // CvTileClickable, ////
    // CvTileExpandable,
    // CvTileSelectable,
    // CvTileStandard,
  },
  props: {
    // expanded: Boolean, ////
    selected: Boolean,
    // NsTile cannot be expandable because of icon
    kind: {
      type: String,
      default: "standard",
      validator: (value) =>
        ["clickable", "selectable", "standard"].includes(value),
    },
    icon: {
      type: [String, Object],
      default: undefined,
      validator(val) {
        if (!val || typeof val === "string") {
          return true;
        }
        return val.render !== null;
      },
    },
  },
  computed: {
    tagType() {
      switch (this.kind) {
        case "clickable":
          return "cv-tile-clickable";
        case "selectable":
          return "cv-tile-selectable";
        default:
          return "cv-tile-standard";
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.tile-icon {
  position: absolute;
  right: 1rem;
  bottom: 1rem;
  flex-shrink: 0;
  width: 1.25rem;
  height: 1.25rem;
}

// .ns-tile:last-child {
//   margin-right: 0;
// }

// .ns-tile {
//   margin-right: $spacing-05;
// }

.pad-bottom {
  padding-bottom: $spacing-10;
}
</style>
