import { NsSystemdServiceCard } from "@nethserver/ns8-ui-lib";
import Cube32 from "@carbon/icons-vue/es/cube/32";

export default {
  title: "Components/NsSystemdServiceCard",
  component: NsSystemdServiceCard,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsSystemdServiceCard },
  template: '<NsSystemdServiceCard v-bind="$props" />',
});

export const Default = Template.bind({});
Default.args = {
  serviceName: "Service name",
  active: true,
  failed: false,
  enabled: true,
  failedLabel: "failed",
  activeLabel: "active",
  inactiveLabel: "inactive",
  enabledLabel: "enabled",
  disaledLabel: "disabled",
  light: true,
  icon: Cube32,
};
