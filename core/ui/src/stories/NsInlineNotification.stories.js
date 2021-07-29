import { NsInlineNotification } from "@nethserver/ns8-ui-lib";

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
  kind: "info",
  showCloseButton: false,
  actionLabel: "",
  closeAriaLabel: "Dismiss notification",
  lowContrast: false,
  loading: false,
};

export const Action = Template.bind({});
Action.args = {
  title: "Something happened",
  description: "Detailed description about what happened",
  actionLabel: "Action button",
  kind: "info",
  showCloseButton: false,
  closeAriaLabel: "Dismiss notification",
  lowContrast: false,
  loading: false,
};

export const Loading = Template.bind({});
Loading.args = {
  title: "Doing stuff",
  description: "Processing data...",
  actionLabel: "Cancel",
  loading: true,
  kind: "info",
  showCloseButton: false,
  closeAriaLabel: "Dismiss notification",
  lowContrast: false,
};
