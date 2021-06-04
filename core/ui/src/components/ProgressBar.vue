<template>
  <!-- indeterminate -->
  <div class="progress-bar-container">
    <div class="slider">
      <template v-if="indeterminate">
        <div class="indeterminate-line"></div>
        <div class="indeterminate-subline inc"></div>
        <div class="indeterminate-subline dec"></div>
      </template>
      <template v-else>
        <div class="line"></div>
        <div class="progress-line" :style="progressLine"></div>
      </template>
    </div>
  </div>
</template>

<script>
export default {
  name: "ProgressBar",
  props: {
    value: {
      // a number between 0 and 100
      type: [Number, String],
      default: 0,
    },
    indeterminate: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      progressLine: {
        width: this.value + "%",
      },
    };
  },
  watch: {
    value: function () {
      this.progressLine.width = this.value + "%";
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/branding";

.progress-bar-container {
  width: 100%;
  position: relative;
}

.slider {
  position: absolute;
  width: 100%;
  height: 5px;
  overflow-x: hidden;
}

.line {
  position: absolute;
  opacity: 0.4;
  background: $interactive-01;
  width: 100%;
  height: 5px;
}

.progress-line {
  position: absolute;
  background: $interactive-01;
  height: 5px;
}

.indeterminate-line {
  position: absolute;
  opacity: 0.4;
  background: $interactive-01;
  width: 150%;
  height: 5px;
}

.indeterminate-subline {
  position: absolute;
  background: $interactive-01;
  height: 5px;
}
.inc {
  animation: increase 2s infinite;
}
.dec {
  animation: decrease 2s 0.5s infinite;
}

@keyframes increase {
  from {
    left: -5%;
    width: 5%;
  }
  to {
    left: 130%;
    width: 100%;
  }
}
@keyframes decrease {
  from {
    left: -80%;
    width: 80%;
  }
  to {
    left: 110%;
    width: 10%;
  }
}
</style>
