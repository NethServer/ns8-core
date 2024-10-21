import { NsTimePicker } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsTimePicker",
  component: NsTimePicker,
  argTypes: {
    dropDirection: {
      control: {
        type: "radio",
        options: ["up", "down", "auto"],
      },
    },
  },
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsTimePicker },
  template: '<NsTimePicker v-bind="$props" />',
});

const TemplateDropDirection = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsTimePicker },
  template: '<NsTimePicker v-bind="$props" style="margin-top: 120px"/>',
});

export const Default = Template.bind({});
Default.args = {
  value: "23:45",
  label: "Label",
  hideClearButton: false,
  invalidMessage: "",
  light: true,
  dropDirection: "down",
};

export const DropUp = TemplateDropDirection.bind({});
DropUp.args = {
  value: "23:45",
  label: "Label",
  hideClearButton: false,
  invalidMessage: "",
  light: true,
  dropDirection: "up",
};
