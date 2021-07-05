<template>
  <div class="empty-state">
    <pictogram title="empty state" class="image">
      <template v-if="hasDefaultSlot">
        <slot></slot>
      </template>
      <template v-else>
        <ExclamationMark />
      </template>
    </pictogram>
    <h5 class="title">{{ title }}</h5>
    <div v-if="description" class="description">
      {{ description }}
    </div>
  </div>
</template>

<script>
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
    description: String,
  },
  computed: {
    hasDefaultSlot() {
      return !!this.$slots.default;
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

.empty-state .title {
  margin-bottom: $spacing-05;
}

.empty-state .description {
  margin-bottom: $spacing-05;
}
</style>
