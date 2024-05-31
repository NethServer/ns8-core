import { NsTextInput } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsTextInput",
  component: NsTextInput,
  argTypes: {
    tooltipAlignment: {
      options: ["start", "center", "end"],
      control: { type: "radio" },
    },
    tooltipDirection: {
      options: ["top", "left", "bottom", "right"],
      control: { type: "radio" },
    },
    type: {
      options: ["text", "password"],
      control: { type: "radio" },
    },
  },
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsTextInput },
  template:
    '<NsTextInput v-bind="$props">\
        <template slot="tooltip">\
          <h6>Tooltip title</h6>\
          <p>Tooltip description</p>\
        </template>\
    </NsTextInput>',
});

export const Default = Template.bind({});
Default.args = {
  label: "Label",
  helperText: "",
  invalidMessage: "",
  warnText: "",
  type: "text",
  passwordHideLabel: "Hide password",
  passwordShowLabel: "Show password",
  passwordVisible: false,
  tooltipAlignment: "start",
  tooltipDirection: "bottom",
  prefix: "",
  light: true,
};
