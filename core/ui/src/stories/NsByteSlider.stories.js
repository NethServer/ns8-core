import { NsByteSlider } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsByteSlider",
  component: NsByteSlider,
  argTypes: {
    tagKind: {
      options: [
        "red",
        "magenta",
        "purple",
        "blue",
        "cyan",
        "teal",
        "green",
        "gray",
        "cool-gray",
        "warm-gray",
        "high-contrast",
      ],
      control: { type: "select" },
    },
    byteUnit: {
      options: ["mib", "gib"],
      control: { type: "radio" },
    },
    onUnlimited: { action: "unlimitedAction" },
    onByteUnit: { action: "byteUnitAction" },
  },
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsByteSlider },
  template:
    '<NsByteSlider v-bind="$props" @unlimited="onUnlimited" @byteUnit="onByteUnit" />',
});

export const Default = Template.bind({});
Default.args = {
  value: "10",
  label: "Storage size",
  min: "1",
  max: "10240",
  step: "1",
  stepMultiplier: "1024",
  minLabel: "",
  maxLabel: "",
  showUnlimited: true,
  unlimitedLabel: "Unlimited",
  limitedLabel: "Specify size",
  isUnlimited: false,
  byteUnit: "gib",
  showHumanReadableLabel: true,
  showMibGibToggle: true,
  invalidMessage: "",
  tagKind: "high-contrast",
  disabled: false,
  light: true,
};
