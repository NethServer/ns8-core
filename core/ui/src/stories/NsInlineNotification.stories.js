import NsInlineNotification from "../components/NsInlineNotification.vue";

export default {
  title: "Components/NsInlineNotification",
  component: NsInlineNotification,
  argTypes: {
    kind: {
      options: ["error", "info", "warning", "success"],
      control: { type: "radio" },
    },
  },
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsInlineNotification },
  template: `<NsInlineNotification v-bind="$props" />`,
});

export const Default = Template.bind({});
Default.args = {
  title: "Something happened",
  description: "Detailed description about what happened",
};

export const Action = Template.bind({});
Action.args = {
  title: "Something happened",
  description: "Detailed description about what happened",
  actionLabel: "Action button",
};
