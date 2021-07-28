import { NsProgressBar } from "andrelib"; ////

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
  value: 75,
  indeterminate: false,
};

export const Indeterminate = Template.bind({});
Indeterminate.args = {
  indeterminate: true,
};
