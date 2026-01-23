import { NsProgress } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsProgress",
  component: NsProgress,
  args: {
    stepId: "firstStep",
    steps: [
      {
        id: "firstStep",
        label: "First Step",
      },
      {
        id: "secondStep",
        label: "Second Step",
      },
      {
        id: "thirdStep",
        label: "Third Step",
      },
      {
        id: "fourthStep",
        label: "Fourth Step",
      },
    ],
  },
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsProgress },
  template: `<NsProgress v-bind="$props" />`,
});

export const Default = Template.bind({});
Default.args = {};

export const SecondStep = Template.bind({});
SecondStep.args = {
  stepId: "secondStep",
};
