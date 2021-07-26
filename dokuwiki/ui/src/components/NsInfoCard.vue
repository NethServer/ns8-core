<template>
  <cv-tile
    kind="standard"
    :light="light"
    :class="['info-card', { 'min-height': minHeight }]"
  >
    <!-- icon -->
    <div v-if="icon" class="row">
      <NsSvg :svg="icon" />
    </div>
    <div class="row">
      <h3 class="title">{{ title }}</h3>
    </div>
    <div class="row">
      <div class="description">{{ description }}</div>
    </div>
  </cv-tile>
</template>

<script>
import NsSvg from "@/components/NsSvg";

export default {
  name: "NsInfoCard",
  components: { NsSvg },
  props: {
    title: {
      type: String,
      required: true,
    },
    description: String,
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
    light: Boolean,
    minHeight: Boolean,
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.info-card {
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 7rem;
}

.min-height {
  // use same value as others Ns*Card components
  min-height: 9rem;
}

.row {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: $spacing-03;
}

.title {
  margin-left: $spacing-02;
  margin-right: $spacing-02;
}

.description {
  margin-left: $spacing-02;
  margin-right: $spacing-02;
}
</style>
