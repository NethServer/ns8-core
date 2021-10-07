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
  cpuUsageWarningTh: 80,
  load1Min: 0.2,
  load5Min: 0.5,
  load15Min: 1.4,
  cpuLoadWarningTh: 1,
  memoryUsage: 88,
  memoryWarningTh: 80,
  swapUsage: 23,
  swapWarningTh: 80,
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
  diskWarningTh: 80,
  light: false,
  slotContent: "Slot content",
};
