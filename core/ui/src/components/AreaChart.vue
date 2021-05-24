<template>
  <ccv-stacked-area-chart
    :data="data"
    :options="options"
  ></ccv-stacked-area-chart>
</template>

<script>
import Vue from "vue";
import "@carbon/charts/styles.css";
import chartsVue from "@carbon/charts-vue";

Vue.use(chartsVue);

export default {
  name: "AreaChart",
  components: {},
  data() {
    return {
      data: [
        {
          group: "Dataset 1",
          date: "2018-12-31T23:00:00.000Z",
          value: 10000,
        },
        {
          group: "Dataset 1",
          date: "2019-01-04T23:00:00.000Z",
          value: 65000,
        },
        {
          group: "Dataset 1",
          date: "2019-01-07T23:00:00.000Z",
          value: 10000,
        },
        {
          group: "Dataset 1",
          date: "2019-01-12T23:00:00.000Z",
          value: 49213,
        },
        {
          group: "Dataset 1",
          date: "2019-01-16T23:00:00.000Z",
          value: 51213,
        },
        {
          group: "Dataset 2",
          date: "2018-12-31T23:00:00.000Z",
          value: 20000,
        },
        {
          group: "Dataset 2",
          date: "2019-01-04T23:00:00.000Z",
          value: 25000,
        },
        {
          group: "Dataset 2",
          date: "2019-01-07T23:00:00.000Z",
          value: 60000,
        },
        {
          group: "Dataset 2",
          date: "2019-01-12T23:00:00.000Z",
          value: 30213,
        },
        {
          group: "Dataset 2",
          date: "2019-01-16T23:00:00.000Z",
          value: 55213,
        },
        {
          group: "Dataset 3",
          date: "2018-12-31T23:00:00.000Z",
          value: 30000,
        },
        {
          group: "Dataset 3",
          date: "2019-01-04T23:00:00.000Z",
          value: 20000,
        },
        {
          group: "Dataset 3",
          date: "2019-01-07T23:00:00.000Z",
          value: 40000,
        },
        {
          group: "Dataset 3",
          date: "2019-01-12T23:00:00.000Z",
          value: 60213,
        },
        {
          group: "Dataset 3",
          date: "2019-01-16T23:00:00.000Z",
          value: 25213,
        },
      ],
      options: {
        title: "Stacked area (time series)",
        axes: {
          left: {
            stacked: true,
            scaleType: "linear",
            mapsTo: "value",
          },
          bottom: {
            scaleType: "time",
            mapsTo: "date",
          },
        },
        curve: "curveMonotoneX",
        height: "400px",
      },
    };
  },
};
</script>
