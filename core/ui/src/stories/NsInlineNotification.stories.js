import { NsInlineNotification } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsInlineNotification",
  component: NsInlineNotification,
  argTypes: {
    kind: {
      options: ["error", "info", "warning", "success"],
      control: { type: "radio" },
    },
    loadingAction: {
      control: { type: "boolean" },
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
  loadingAction: false,
  closeAriaLabel: "Dismiss notification",
  lowContrast: true,
  loading: false,
  timer: null,
};

const TitleSlotTemplate = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsInlineNotification },
  template: `<NsInlineNotification v-bind="$props">
                <template slot="title">
                  Title slot
                </template>
            </NsInlineNotification>`,
});

export const TitleSlot = TitleSlotTemplate.bind({});
TitleSlot.args = {
  title: "",
  description: "Detailed description about what happened",
  kind: "info",
  showCloseButton: false,
  actionLabel: "",
  loadingAction: false,
  closeAriaLabel: "Dismiss notification",
  lowContrast: true,
  loading: false,
  timer: null,
};

const DescriptionSlotTemplate = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsInlineNotification },
  template: `<NsInlineNotification v-bind="$props">
                <template slot="description">
                  Description slot
                </template>
            </NsInlineNotification>`,
});

export const DescriptionSlot = DescriptionSlotTemplate.bind({});
DescriptionSlot.args = {
  title: "Something happened",
  description: "",
  kind: "info",
  showCloseButton: false,
  actionLabel: "",
  loadingAction: false,
  closeAriaLabel: "Dismiss notification",
  lowContrast: true,
  loading: false,
  timer: null,
};

export const Action = Template.bind({});
Action.args = {
  title: "Something happened",
  description: "Detailed description about what happened",
  actionLabel: "Action button",
  kind: "info",
  showCloseButton: false,
  loadingAction: false,
  closeAriaLabel: "Dismiss notification",
  lowContrast: true,
  loading: false,
  timer: null,
};

export const Loading = Template.bind({});
Loading.args = {
  title: "Doing stuff",
  description: "Processing data...",
  actionLabel: "Cancel",
  loading: true,
  kind: "info",
  showCloseButton: false,
  loadingAction: false,
  closeAriaLabel: "Dismiss notification",
  lowContrast: true,
  timer: null,
};

export const LoadingAction = Template.bind({});
LoadingAction.args = {
  title: "Preparing download",
  description: "Please wait while the backup is generated...",
  actionLabel: "Download backup",
  loading: false,
  loadingAction: true,
  kind: "info",
  showCloseButton: false,
  closeAriaLabel: "Dismiss notification",
  lowContrast: true,
  timer: null,
};

export const Timer = Template.bind({});
Timer.args = {
  title: "Hey there",
  description: "Something will happen in a moment...",
  actionLabel: "Cancel",
  loading: false,
  loadingAction: false,
  kind: "info",
  showCloseButton: false,
  closeAriaLabel: "Dismiss notification",
  lowContrast: true,
  timer: 5000,
};
