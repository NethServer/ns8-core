<template>
  <ccv-meter-chart :data="data" :options="options"></ccv-meter-chart>
</template>

<script>
import Vue from "vue";
import "@carbon/charts/styles.css";
import chartsVue from "@carbon/charts-vue";

// IBM Plex should either be imported in your project by using Carbon
// or consumed manually through an import
// import "../plex-and-carbon-components.css"; //// needed?

Vue.use(chartsVue);

//// move to ns8-ui-lib

export default {
  name: "MeterChart",
  components: {},
  props: {
    title: { type: String, default: "" },
    label: { type: String, default: "" },
    value: { type: Number, default: 0 },
    warningTh: { type: Number, default: 70 },
    dangerTh: { type: Number, default: 90 },
    height: { type: String, default: "3rem" },
  },
  data() {
    return {
      data: [],
      options: {},
    };
  },
  watch: {
    ////
    value: function () {
      console.log("watch value", this.value); ////
      this.data[0].value = this.value;
    },
  },
  created() {
    this.data = [
      {
        group: this.label,
        value: this.value,
      },
    ];

    this.options = {
      title: this.title,
      meter: {
        peak: null,
        status: {
          ranges: [
            {
              range: [0, this.warningTh],
              status: "success",
            },
            {
              range: [this.warningTh, this.dangerTh],
              status: "warning",
            },
            {
              range: [this.dangerTh, 100],
              status: "danger",
            },
          ],
        },
      },
      height: this.height,
    };
  },
};
</script>
