import { NsSystemLogsCard } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsSystemLogsCard",
  component: NsSystemLogsCard,
  argTypes: {
    context: {
      options: ["cluster", "node", "module"],
      control: { type: "radio" },
    },
  },
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsSystemLogsCard },
  template: '<NsSystemLogsCard v-bind="$props" />',
});

export const Default = Template.bind({});
Default.args = {
  title: "Logs",
  description: "Show nextcloud1 logs",
  buttonLabel: "Go to System logs",
  light: true,
  router: null,
  context: "module",
  nodeId: "",
  moduleId: "nextcloud1",
  searchQuery: "",
  maxLines: "500",
  startDate: "2022-01-01",
  startTime: "00:00",
  endDate: "2022-12-31",
  endTime: "23:59",
  followLogs: false,
};
