import { NsPasswordInput } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsPasswordInput",
  component: NsPasswordInput,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsPasswordInput },
  template: '<NsPasswordInput v-bind="$props" />',
});

export const Default = Template.bind({});
Default.args = {
  value: "",
  label: "Password",
  invalidMessage: "",
  helperText: "",
  passwordHideLabel: "Hide password",
  passwordShowLabel: "Show password",
  passwordVisible: false,
  minLength: 8,
  lengthLabel: "Long enough",
  lowercaseLabel: "Lowercase letter",
  uppercaseLabel: "Uppercase letter",
  numberLabel: "Number",
  symbolLabel: "Symbol",
  light: false,
};
