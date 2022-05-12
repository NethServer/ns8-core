import { NsMeterChart } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsMeterChart",
  component: NsMeterChart,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsMeterChart },
  template: '<NsMeterChart v-bind="$props" />',
});

export const Default = Template.bind({});
Default.args = {
  label: "Usage",
  value: 65,
  loading: false,
  warningThreshold: 70,
  dangerThreshold: 90,
  progressBarHeight: "10px",
};
