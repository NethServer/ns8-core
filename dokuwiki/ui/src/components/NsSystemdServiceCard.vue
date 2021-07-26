<template>
  <cv-tile
    kind="standard"
    :light="light"
    :class="['service-card', { 'min-height': minHeight }]"
  >
    <!-- icon -->
    <div v-if="icon" class="row">
      <NsSvg :svg="icon" />
    </div>
    <div class="row">
      <h3 class="service-name">{{ serviceName }}</h3>
    </div>
    <div class="row">
      <div class="stats">
        <cv-tag v-if="failed" kind="red" :label="failedLabel"></cv-tag>
        <cv-tag v-if="active" kind="green" :label="activeLabel"></cv-tag>
        <cv-tag v-else kind="gray" :label="inactiveLabel"></cv-tag>
        <cv-tag v-if="enabled" kind="blue" :label="enabledLabel"></cv-tag>
        <cv-tag v-else kind="gray" :label="disaledLabel"></cv-tag>
      </div>
    </div>
  </cv-tile>
</template>

<script>
import NsSvg from "@/components/NsSvg";

export default {
  name: "NsSystemdServiceCard",
  components: { NsSvg },
  props: {
    serviceName: {
      type: String,
      required: true,
    },
    active: Boolean,
    failed: Boolean,
    enabled: Boolean,
    failedLabel: {
      type: String,
      default: "failed",
    },
    activeLabel: {
      type: String,
      default: "active",
    },
    inactiveLabel: {
      type: String,
      default: "inactive",
    },
    enabledLabel: {
      type: String,
      default: "enabled",
    },
    disaledLabel: {
      type: String,
      default: "disabled",
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
    light: Boolean,
    minHeight: Boolean,
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.service-card {
  display: flex;
  flex-direction: column;
  justify-content: center;
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

.service-name {
  margin-left: $spacing-02;
  margin-right: $spacing-02;
}

.success-icon {
  color: $support-02;
  margin-right: $spacing-02;
}

.error-icon {
  color: $danger-01;
  margin-right: $spacing-02;
}

.warning-icon {
  color: $support-03;
  margin-right: $spacing-02;
}

.stats {
  display: flex;
  align-items: center;
  margin-left: $spacing-03;
  margin-right: $spacing-03;
}
</style>
