import { NsNodeCard } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsNodeCard",
  component: NsNodeCard,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsNodeCard },
  template: '<NsNodeCard v-bind="$props">{{ slotContent }}</NsNodeCard>',
});

export const Default = Template.bind({});
Default.args = {
  nodeId: "1",
  nodeLabel: "Node",
  isLeader: true,
  leaderLabel: "leader",
  workerLabel: "worker",
  cpuUsageLabel: "CPU usage",
  cpuLoadLabel: "CPU load",
  cpuLoadTooltip: "CPU average load of last 1 / 5 / 15 minutes",
  memoryUsageLabel: "Memory usage",
  swapUsageLabel: "Swap usage",
  diskUsageLabel: "usage",
  cpuUsage: 42,
  cpuUsageWarningTh: 90,
  load1Min: 0.2,
  load5Min: 0.5,
  load15Min: 1.4,
  cpuLoadWarningTh: 90,
  memoryUsage: 88,
  memoryWarningTh: 90,
  swapUsage: 23,
  swapWarningTh: 90,
  disksUsage: [
    {
      diskId: "Disk 1",
      usage: 55,
    },
    {
      diskId: "Disk 2",
      usage: 93,
    },
  ],
  diskWarningTh: 90,
  light: false,
  slotContent: "Slot content",
};
