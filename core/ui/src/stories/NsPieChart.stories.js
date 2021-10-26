import { NsPieChart } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsPieChart",
  component: NsPieChart,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsPieChart },
  template: '<NsPieChart v-bind="$props" />',
});

export const Default = Template.bind({});
Default.args = {
  loading: false,
  title: "Title",
  height: "18rem",
  data: [
    {
      group: "Group 1",
      value: 20000,
    },
    {
      group: "Group 2",
      value: 65000,
    },
    {
      group: "Group 3",
      value: 75000,
    },
    {
      group: "Group 4",
      value: 1200,
    },
    {
      group: "Group 5",
      value: 10000,
    },
    {
      group: "Misc",
      value: 25000,
    },
  ],
};
