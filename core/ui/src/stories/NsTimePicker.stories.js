import { NsTimePicker } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsTimePicker",
  component: NsTimePicker,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsTimePicker },
  template: '<NsTimePicker v-bind="$props" />',
});

export const Default = Template.bind({});
Default.args = {
  value: "23:45",
  label: "Label",
  hideClearButton: false,
  invalidMessage: "",
  light: true,
};

export const Drop_Up = Template.bind({});
Drop_Up.args = {
  value: "23:45",
  label: "Label",
  hideClearButton: false,
  invalidMessage: "",
  light: true,
};
