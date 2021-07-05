<template>
  <div
    data-notification
    :class="[
      'cv-inline-notification',
      `${carbonPrefix}--inline-notification`,
      `${carbonPrefix}--inline-notification--${kind.toLowerCase()}`,
      {
        [`${carbonPrefix}--inline-notification--low-contrast`]: lowContrast,
      },
    ]"
    v-on="$listeners"
    :role="isAlert ? 'alert' : false"
    :aria-live="!isAlert ? 'polite' : false"
  >
    <div :class="`${carbonPrefix}--inline-notification__details details`">
      <component
        :is="icon"
        :class="`${carbonPrefix}--inline-notification__icon`"
      />
      <div :class="`${carbonPrefix}--inline-notification__text-wrapper`">
        <p
          :class="`${carbonPrefix}--inline-notification__title title`"
          v-html="title"
        ></p>
        <p
          v-if="description"
          :class="[
            `${carbonPrefix}--inline-notification__subtitle`,
            { 'mg-right': loading },
          ]"
          v-html="description"
        ></p>
        <p v-if="loading" class="loader"></p>
      </div>
    </div>
    <button
      v-if="actionLabel"
      @click="$emit('action')"
      :class="[
        `${carbonPrefix}--inline-notification__action-button`,
        `${carbonPrefix}--btn`,
        `${carbonPrefix}--btn--sm`,
        `${carbonPrefix}--btn--ghost`,
      ]"
      type="button"
    >
      {{ actionLabel }}
    </button>
    <button
      v-if="showCloseButton"
      type="button"
      :aria-label="closeAriaLabel"
      data-notification-btn
      :class="`${carbonPrefix}--inline-notification__close-button`"
      @click="$emit('close')"
    >
      <Close20 :class="`${carbonPrefix}--inline-notification__close-icon`" />
    </button>
  </div>
</template>

<script>
import { CvInlineNotification } from "@carbon/vue";

export default {
  name: "NsInlineNotification",
  extends: CvInlineNotification,
  props: {
    showCloseButton: {
      type: Boolean,
      default: false,
    },
    description: String,
    actionLabel: { type: String, default: "" },
    closeAriaLabel: { type: String, default: "Dismiss notification" },
    kind: {
      type: String,
      default: "info",
      validator: (val) => ["error", "info", "warning", "success"].includes(val),
    },
    lowContrast: Boolean,
    loading: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.title {
  margin-right: 0.75rem;
  margin-bottom: 0.2rem;
}

.mg-right {
  margin-right: 0.75rem;
}

// place action button next to description
.details {
  flex-grow: 0;
}

// place close button on the right
.bx--inline-notification__close-button {
  position: absolute !important;
}
</style>
