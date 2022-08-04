import { NsSlider } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsSlider",
  component: NsSlider,
  argTypes: {
    onUnlimited: { action: "unlimitedAction" },
  },
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsSlider },
  template: '<NsSlider v-bind="$props" @unlimited="onUnlimited" />',
});

export const Default = Template.bind({});
Default.args = {
  value: "10",
  label: "Holiday duration",
  unitLabel: "Days",
  min: "1",
  max: "365",
  step: "1",
  stepMultiplier: "10",
  minLabel: "",
  maxLabel: "",
  showUnlimited: true,
  unlimitedLabel: "Forever",
  limitedLabel: "Specify duration",
  isUnlimited: false,
  invalidMessage: "",
  disabled: false,
  light: true,
};
