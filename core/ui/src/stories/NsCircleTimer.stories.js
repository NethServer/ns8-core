import { NsCircleTimer } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsCircleTimer",
  component: NsCircleTimer,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsCircleTimer },
  template: '<NsCircleTimer v-bind="$props" />',
});

export const Default = Template.bind({});
Default.args = {
  timeLimit: 5,
  size: 16,
  strokeWidth: 20,
  color: "#00a2de",
};
