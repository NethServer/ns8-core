<template>
  <div
    class="logs-output cv-code-snippet-multiline"
    :class="[
      `bx--snippet`,
      `bx--snippet--multi`,
      {
        'bx--snippet--disabled': disabled,
        'bx--snippet--wraptext': wrapText,
        'bx--snippet--light': isLight,
      },
    ]"
    data-code-snippet
  >
    <!-- <div ////
      @scroll="handleScroll"
      class="bx--snippet-container"
      ref="logsContainer"
    > -->
    <cv-skeleton-text
      v-if="loading"
      :paragraph="true"
      :line-count="4"
    ></cv-skeleton-text>
    <div v-else class="bx--snippet-container" ref="logsContainer">
      <pre><slot /></pre>
    </div>
  </div>
</template>

<script>
import { carbonPrefixMixin, themeMixin } from "@carbon/vue/src/mixins";

export default {
  name: "LogsOutput",
  props: {
    id: {
      type: String,
      required: true,
    },
    scrollToBottom: {
      type: Boolean,
      default: true,
    },
    loading: Boolean,
    wrapText: Boolean,
    light: Boolean,
    disabled: Boolean,
  },
  mixins: [carbonPrefixMixin, themeMixin],
  watch: {
    id: function () {
      //// handle reqId
      this.$root.$on("logsStart", this.logsUpdated);
    },
  },
  created() {
    //// handle reqId
    this.$root.$on("logsStart", this.logsUpdated);
  },
  beforeDestroy() {
    //// handle reqId
    this.$root.$off("logsStart");
  },
  methods: {
    // handleScroll: function (el) { ////
    //   // this method is invoked only on user scrolls
    //   if (
    //     el.target.offsetHeight + el.target.scrollTop >=
    //     el.target.scrollHeight
    //   ) {
    //     this.scrolledToBottom = true;
    //   } else {
    //     this.scrolledToBottom = false;
    //   }

    //   // if (!this.userHasScrolled) { ////
    //   //   this.userHasScrolled = true;
    //   // }
    // },
    logsUpdated() {
      // new log lines are displayed

      setTimeout(() => {
        // keep scrolling to the bottom if user has never scrolled or if he has scrolled to the bottom
        if (this.scrollToBottom) {
          this.scrollToBottomOfContainer();
        }
      }, 100);
    },
    scrollToBottomOfContainer() {
      const el = this.$refs.logsContainer;
      el.scrollTop = el.scrollHeight;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";
</style>

<style lang="scss">
@import "../styles/carbon-utils";

// global styles

.bx--snippet--multi .bx--snippet-container pre {
  overflow-x: visible;
}
</style>
