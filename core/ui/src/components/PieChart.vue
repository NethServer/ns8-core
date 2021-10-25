<template>
  <ccv-pie-chart :data="data" :options="options"></ccv-pie-chart>
</template>

<script>
import Vue from "vue";
import "@carbon/charts/styles.css";
import chartsVue from "@carbon/charts-vue";

Vue.use(chartsVue);

//// move to ns8-ui-lib

export default {
  name: "PieChart",
  components: {},
  props: {
    data: Array,
    loading: { type: Boolean, default: false },
    title: { type: String, default: "" },
    height: { type: String, default: "15rem" },
  },
  data() {
    return {
      options: {},
    };
  },
  watch: {
    loading: function () {
      if (this.loading) {
        // show skeleton
        this.options.data.loading = true;
      } else {
        this.options.data.loading = false;
      }
    },
  },
  created() {
    this.options = {
      title: this.title,
      resizable: true,
      // legend: { ////
      //   alignment: "center",
      // },
      // pie: {
      //   alignment: "center",
      // },
      data: { loading: this.loading },
      height: this.height,
    };
  },
};
</script>
