import { NsToastNotification } from "andrelib"; ////

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
  showCloseButton: true,
  actionLabel: "Action",
  timestamp: new Date(),
};

export const Progress = Template.bind({});
Progress.args = {
  title: "Task in progress",
  description: "Current step description",
  showCloseButton: true,
  actionLabel: "Action",
  progress: 75,
  isProgressShown: true,
  timestamp: new Date(),
};
