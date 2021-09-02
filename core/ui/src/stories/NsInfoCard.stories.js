import { NsInfoCard } from "@nethserver/ns8-ui-lib";
import Application32 from "@carbon/icons-vue/es/application/32";

export default {
  title: "Components/NsInfoCard",
  component: NsInfoCard,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsInfoCard },
  template: '<NsInfoCard v-bind="$props" />',
});

export const Default = Template.bind({});
Default.args = {
  title: "Card title",
  description: "Card description",
  light: false,
  icon: Application32,
};
