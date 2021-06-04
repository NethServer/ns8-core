import ProgressBar from "../components/ProgressBar.vue";

export default {
  title: "Components/ProgressBar",
  component: ProgressBar,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { ProgressBar },
  template: '<ProgressBar v-bind="$props" />',
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
