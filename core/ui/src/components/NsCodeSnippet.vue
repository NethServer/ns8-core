<template>
  <!--  //// rename to NsCodeSnippet -->
  <div class="cv-code-snippet-multiline" :class="classes" data-code-snippet>
    <div
      :class="[
        `${carbonPrefix}--snippet-container`,
        { 'no-expand': hideExpandButton },
      ]"
    >
      <code ref="code">
        <pre><slot /></pre>
      </code>
      <!-- textarea cannot be fully hidden for clipboard to work -->
      <textarea
        class="cv-code-snippet__clippy"
        style="
          position: absolute;
          left: -9999px;
          max-width: 0;
          opacity: 0;
          overflow: hidden;
        "
        aria-hidden="true"
        ref="clippy"
      ></textarea>
    </div>
    <cv-icon-button
      v-if="!hideCopyButton"
      :class="`${carbonPrefix}--copy-btn`"
      kind="ghost"
      :icon="Copy20"
      :label="tooltipToShow"
      tip-position="bottom"
      tip-alignment="center"
      @click="copyCode"
      :disabled="disabled"
    />
    <cv-button
      v-if="!hideExpandButton"
      type="button"
      kind="ghost"
      size="small"
      :class="`${carbonPrefix}--snippet-btn--expand`"
      @click="toggleExpand"
    >
      <span :class="`${carbonPrefix}--snippet-btn--text`">{{
        expandButtonText
      }}</span>
      <ChevronDown20 :class="`${carbonPrefix}--icon-chevron--down`" />
    </cv-button>
  </div>
</template>

<script>
import { CvButton, CvIconButton } from "@carbon/vue";
import Copy20 from "@carbon/icons-vue/es/copy/20";
import ChevronDown20 from "@carbon/icons-vue/es/chevron--down/20";
import { carbonPrefixMixin, themeMixin } from "@carbon/vue/src/mixins";

export default {
  name: "NsCodeSnippet",
  components: { CvButton, CvIconButton, ChevronDown20 },
  mixins: [carbonPrefixMixin, themeMixin],
  props: {
    copyTooltip: {
      type: String,
      default: "Copy to clipboard",
    },
    copyFeedback: { type: String, default: "Copied" },
    wrapText: { type: Boolean, default: true },
    disabled: Boolean,
    feedbackAriaLabel: String,
    hideCopyButton: Boolean,
    lessText: { type: String, default: "Show less" },
    moreText: { type: String, default: "Show more" },
    hideExpandButton: Boolean,
  },
  data() {
    return {
      expanded: false,
      Copy20, //// use mixin?
      justCopied: false,
    };
  },
  computed: {
    classes() {
      return [
        `cv-code-snippet`,
        `${this.carbonPrefix}--snippet`,
        `${this.carbonPrefix}--snippet--multi`,
        {
          [`${this.carbonPrefix}--snippet--light`]: this.isLight,
          [`${this.carbonPrefix}--snippet--no-copy`]: this.hideCopyButton,
          [`${this.carbonPrefix}--snippet--wraptext`]: this.wrapText,
          [`${this.carbonPrefix}--snippet--expand`]: this.expanded,
          [`${this.carbonPrefix}--snippet--disabled`]: this.disabled,
        },
      ];
    },
    tooltipToShow() {
      if (this.justCopied) {
        return this.copyFeedback;
      } else {
        return this.copyTooltip;
      }
    },
    expandButtonText() {
      return this.expanded ? this.lessText : this.moreText;
    },
  },
  methods: {
    toggleExpand() {
      this.expanded = !this.expanded;
    },
    copyCode() {
      // copy to clipboard
      this.$refs.clippy.value = this.$refs.code.innerText;
      this.$refs.clippy.select();
      document.execCommand("copy");

      this.justCopied = true;
      setTimeout(() => {
        this.justCopied = false;
      }, 2000);
    },
  },
};
</script>

<style scoped lang="scss">
.bx--snippet--multi .bx--snippet-container {
  min-height: 4.5rem;
  max-height: 9.6rem;
  overflow-y: hidden;
}

.bx--snippet--multi .bx--snippet-container.no-expand {
  min-height: 2.5rem;
}

.bx--snippet--multi.bx--snippet--expand .bx--snippet-container {
  max-height: 100%;
}

.bx--snippet--multi .bx--copy-btn {
  position: absolute !important;
  top: 0.5rem;
  right: 0.5rem;
  z-index: 11;
  width: 3rem;
  height: 2rem;
}

.bx--snippet-btn--expand {
  right: 0.5rem;
  bottom: 0.5rem;
}
</style>
