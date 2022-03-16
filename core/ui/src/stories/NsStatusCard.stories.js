import { NsStatusCard } from "@nethserver/ns8-ui-lib";
import Cube32 from "@carbon/icons-vue/es/cube/32";

export default {
  title: "Components/NsStatusCard",
  component: NsStatusCard,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsStatusCard },
  template: '<NsStatusCard v-bind="$props" />',
});

export const Default = Template.bind({});
Default.args = {
  value: 5,
  valueSuccess: 3,
  valueError: 2,
  valueWarning: 1,
  label: "Components",
  errorTooltip: "Failed components",
  successTooltip: "Running components",
  warningTooltip: "Unknown components",
  light: true,
  icon: Cube32,
};
