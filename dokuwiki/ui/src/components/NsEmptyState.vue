<template>
  <div class="empty-state">
    <pictogram title="empty state" class="image">
      <template v-if="hasPictogramSlot">
        <slot name="pictogram"></slot>
      </template>
      <template v-else>
        <ExclamationMark />
      </template>
    </pictogram>
    <h5 class="title">{{ title }}</h5>

    <template v-if="hasDescriptionSlot">
      <div class="description">
        <slot name="description"></slot>
      </div>
    </template>
  </div>
</template>

<script>
//// FILE COPIED FROM CORE

import Pictogram from "./Pictogram.vue";
import ExclamationMark from "./pictograms/ExclamationMark";

export default {
  name: "NsEmptyState",
  components: { Pictogram, ExclamationMark },
  props: {
    title: {
      type: String,
      required: true,
    },
  },
  computed: {
    hasPictogramSlot() {
      return !!this.$slots.pictogram;
    },
    hasDescriptionSlot() {
      return !!this.$slots.description;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.empty-state {
  text-align: center;
  padding: $spacing-07;
}

.empty-state .image {
  margin-bottom: $spacing-05;
}

.empty-state .description {
  margin-top: $spacing-05;
}
</style>
