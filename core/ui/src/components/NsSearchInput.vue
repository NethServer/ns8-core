<template>
  <div :class="`cv-text-input ${carbonPrefix}--form-item`">
    <label
      :for="uid"
      :class="[
        `${carbonPrefix}--label`,
        {
          [`${carbonPrefix}--label--disabled`]:
            $attrs.disabled !== undefined && $attrs.disabled,
        },
      ]"
      >{{ label }}</label
    >
    <div
      :class="`${carbonPrefix}--text-input__field-wrapper`"
      :data-invalid="isInvalid"
    >
      <Search16 :class="`search-icon`" @click="focusSearchInput" />
      <!-- //// i18n label -->
      <cv-icon-button
        kind="ghost"
        :icon="Close16"
        :label="clearLabel"
        v-if="!isInvalid"
        :class="`clear-icon`"
        @click="clearSearchInput"
        tip-position="bottom"
        tip-alignment="end"
      />
      <!-- <Close16 ////
        v-if="query && !isInvalid"
        :class="`clear-icon`"
        @click="clearSearchInput"
      /> -->
      <WarningFilled16
        v-if="isInvalid"
        :class="`${carbonPrefix}--text-input__invalid-icon`"
      />
      <input
        :id="uid"
        :class="[
          `${carbonPrefix}--text-input search-input`,
          {
            [`${carbonPrefix}--text-input--light`]: isLight,
            [`${carbonPrefix}--text-input--invalid`]: isInvalid,
          },
        ]"
        v-bind="$attrs"
        :value="value"
        v-on="inputListeners"
        ref="input"
        v-debounce="search"
      />
    </div>
    <div :class="`${carbonPrefix}--form-requirement`" v-if="isInvalid">
      <slot name="invalid-message">{{ invalidMessage }}</slot>
    </div>
    <div
      v-if="!isInvalid && isHelper"
      :class="[
        `${carbonPrefix}--form__helper-text`,
        { [`${carbonPrefix}--form__helper-text--disabled`]: $attrs.disabled },
      ]"
    >
      <slot name="helper-text">{{ helperText }}</slot>
    </div>
  </div>
</template>

<script>
import { CvTextInput, CvIconButton } from "@carbon/vue";
import Search16 from "@carbon/icons-vue/es/search/16";
import Close16 from "@carbon/icons-vue/es/close/16";

export default {
  name: "NsSearchInput",
  extends: CvTextInput,
  components: { Search16, CvIconButton },
  props: {
    helperText: { type: String, default: undefined },
    invalidMessage: { type: String, default: undefined },
    label: String,
    value: String,
    clearLabel: { type: String, default: "Clear search" },
  },
  data() {
    return {
      query: "",
      Close16,
    };
  },
  // computed: {
  //   query() {
  //     if (this.$refs.input) {
  //       return this.$refs.input;
  //     } else {
  //       return "";
  //     }
  //   },
  //   isQueryTyped() {
  //     if (this.$refs.input) {
  //       return this.$refs.input.value.length > 0;
  //     } else {
  //       console.log("ref not found", this.$refs); ////
  //       return false;
  //     }
  //   },
  // },
  // updated() { ////
  //   console.log("updated", this.$refs.input); ////
  //   if (!this.$refs.input) {
  //     return;
  //   }

  //   this.query = this.$refs.input.value;
  // },
  methods: {
    // onInput(ev) { ////
    //   console.log("onInput", ev); ////
    //   this.query = this.$refs.input.value;
    // },
    focusSearchInput() {
      const searchInput = this.$refs.input;
      searchInput.focus();
    },
    search() {
      this.query = this.$refs.input.value;

      console.log("search!", this.query); ////

      this.$emit("search", this.query);
    },
    // isQueryTyped() {
    //   ////
    //   console.log("isQueryTyped", this.$refs.input); ////

    //   if (this.$refs.input) {
    //     return this.$refs.input.value.length > 0;
    //   } else {
    //     console.log("ref not found", this.$refs); ////
    //     return false;
    //   }
    // },
    clearSearchInput() {
      this.$emit("clear");

      // this.$refs.input.value = ""; ////
      // const context = this;
      // setTimeout(() => {
      //   context.query = "";
      // }, 550);
    },
  },
};
</script>

<style scoped lang="scss">
.search-icon {
  position: absolute;
  top: 50%;
  left: 1rem;
  -webkit-transform: translateY(-50%);
  transform: translateY(-50%);
}

.search-input {
  padding-left: 3rem;
}

.clear-icon {
  position: absolute !important;
  top: 50%;
  right: 0.5rem;
  -webkit-transform: translateY(-50%);
  transform: translateY(-50%);
  padding: 0.5rem;
  min-height: 0;
  // width: 2rem;
  // height: 2rem;
  //// width, height, cursor
}
</style>
