import { NsCodeSnippet } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsCodeSnippet",
  component: NsCodeSnippet,
  argTypes: {
    copyTipPosition: {
      options: ["top", "right", "bottom", "left"],
      control: { type: "radio" },
    },
    copyTipAlignment: {
      options: ["start", "center", "end"],
      control: { type: "radio" },
    },
  },
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsCodeSnippet },
  template: `<NsCodeSnippet v-bind="$props">{{ slotContent }}</NsCodeSnippet>`,
});

export const Default = Template.bind({});
Default.args = {
  wrapText: true,
  expanded: false,
  hideCopyButton: false,
  lessText: "Show less",
  moreText: "Show more",
  hideExpandButton: false,
  copyTooltip: "Copy to clipboard",
  copyFeedback: "Copied!",
  feedbackAriaLabel: "Copied!",
  copyTipPosition: "bottom",
  copyTipAlignment: "end",
  disabled: false,
  light: false,
  slotContent: `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. 
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo.
Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt.
Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt ut labore et dolore magnam aliquam quaerat voluptatem.
Ut enim ad minima veniam, quis nostrum exercitationem ullam corporis suscipit laboriosam, nisi ut aliquid ex ea commodi consequatur?
Quis autem vel eum iure reprehenderit qui in ea voluptate velit esse quam nihil molestiae consequatur, vel illum qui dolorem eum fugiat quo voluptas nulla pariatur?`,
};

export const OneLine = Template.bind({});
OneLine.args = {
  hideExpandButton: true,
  wrapText: true,
  expanded: false,
  hideCopyButton: false,
  lessText: "Show less",
  moreText: "Show more",
  copyTooltip: "Copy to clipboard",
  copyFeedback: "Copied",
  copyTipPosition: "bottom",
  copyTipAlignment: "end",
  disabled: false,
  light: false,
  slotContent: `let foo = "bar";`,
};
