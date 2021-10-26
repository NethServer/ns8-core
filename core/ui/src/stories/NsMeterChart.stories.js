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
  title: "Title",
  label: "Label",
  value: 65,
  warningTh: 70,
  dangerTh: 90,
  height: "6rem",
};
