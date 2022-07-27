import { NsProgressBar } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsProgressBar",
  component: NsProgressBar,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsProgressBar },
  template: '<NsProgressBar v-bind="$props" />',
});

export const Determinate = Template.bind({});
Determinate.args = {
  value: 65,
  indeterminate: false,
  height: "5px",
  warningThreshold: 70,
  dangerThreshold: 90,
  useStatusColors: false,
  useHealthyColor: true,
};

export const Indeterminate = Template.bind({});
Indeterminate.args = {
  value: 65,
  indeterminate: true,
  height: "5px",
  warningThreshold: 70,
  dangerThreshold: 90,
  useStatusColors: false,
  useHealthyColor: true,
};
