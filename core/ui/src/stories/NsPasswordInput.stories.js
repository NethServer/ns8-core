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
  newPaswordHelperText: "",
  confirmPaswordHelperText: "",
  newPasswordInvalidMessage: "",
  confirmPasswordInvalidMessage: "",
  newPasswordLabel: "New password",
  confirmPasswordLabel: "Re-enter new password",
  passwordHideLabel: "Hide password",
  passwordShowLabel: "Show password",
  newPasswordVisible: false,
  confirmPasswordVisible: false,
  focus: {},
  light: true,
  minLength: 8,
  checkComplexity: true,
  clearConfirmPasswordCommand: 0,
  lengthLabel: "Long enough",
  lowercaseLabel: "Lowercase letter",
  uppercaseLabel: "Uppercase letter",
  numberLabel: "Number",
  symbolLabel: "Symbol",
  equalLabel: "Equal",
  disabled: false,
};
