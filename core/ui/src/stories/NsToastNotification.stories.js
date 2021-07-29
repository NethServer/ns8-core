import { NsToastNotification } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsToastNotification",
  component: NsToastNotification,
  argTypes: {
    kind: {
      options: ["error", "info", "warning", "success"],
      control: { type: "radio" },
    },
  },
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsToastNotification },
  template: '<NsToastNotification v-bind="$props" />',
});

export const Default = Template.bind({});
Default.args = {
  title: "Something happened",
  description: "Detailed description about what happened",
  kind: "info",
  showCloseButton: true,
  actionLabel: "Action",
  lowContrast: false,
  timestamp: new Date(),
  read: false,
  progress: 0,
  isProgressShown: false,
  closeAriaLabel: "Dismiss notification",
};

export const Progress = Template.bind({});
Progress.args = {
  title: "Task in progress",
  description: "Current step description",
  kind: "info",
  showCloseButton: true,
  actionLabel: "Action",
  lowContrast: false,
  progress: 75,
  isProgressShown: true,
  timestamp: new Date(),
  read: false,
  closeAriaLabel: "Dismiss notification",
};
