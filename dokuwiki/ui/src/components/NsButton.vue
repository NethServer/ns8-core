<template>
  <button
    class="cv-button"
    :class="buttonClassOpts()"
    v-on="inputListeners"
    role="button"
  >
    <slot></slot>

    <!-- spinning loader -->
    <span
      v-if="loading"
      :class="['loader', `${carbonPrefix}--btn__icon`]"
    ></span>

    <!-- icon -->
    <CvSvg
      v-if="icon && !loading"
      :svg="icon"
      :class="`${carbonPrefix}--btn__icon`"
    />
  </button>
</template>

<script>
//// FILE COPIED FROM CORE

import { CvButton } from "@carbon/vue";

export default {
  name: "NsButton",
  extends: CvButton,
  props: {
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
    kind: {
      type: String,
      default: "default",
      validator: (val) =>
        [
          "default",
          "primary",
          "secondary",
          "tertiary",
          "ghost",
          "danger",
          "danger--ghost",
          "danger--tertiary",
        ].includes(val),
    },
    size: {
      type: String,
      default: "default",
      validator: (val) =>
        ["default", "field", "small", "sm", "lg", "xl"].includes(val),
    },
    loading: Boolean,
  },
};
</script>

<style scoped lang="scss">
// icon size
.bx--btn .bx--btn__icon {
  width: 1.25rem;
  height: 1.25rem;
}

// icon size for small button
.bx--btn.bx--btn--sm .bx--btn__icon {
  width: 1rem;
  height: 1rem;
}

.loader {
  right: 1rem;
  border: 3px solid transparent;
  border-radius: 50%;
  border-top: 3px solid currentColor;
  border-right: 3px solid currentColor;
  border-bottom: 3px solid currentColor;
  animation: spin 0.5s linear infinite;
}

// fix loader position for ghost button
.bx--btn--ghost .bx--btn__icon.loading {
  position: static;
  margin-left: 0.5rem;
}

// fix loader position for ghost danger button
.bx--btn--danger-ghost .bx--btn__icon.loading,
.bx--btn--danger--ghost .bx--btn__icon.loading {
  position: static;
  margin-left: 0.5rem;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>
