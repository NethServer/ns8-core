import { NsComboSearchBox } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsComboSearchBox",
  component: NsComboSearchBox,
  argTypes: {
    tooltipAlignment: {
      options: ["start", "center", "end"],
      control: { type: "radio" },
    },
    tooltipDirection: {
      options: ["top", "left", "bottom", "right"],
      control: { type: "radio" },
    },
    testSearch: { action: "searchAction" },
    testChange: { action: "changeAction" },
    testFilter: { action: "filterAction" },
  },
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsComboSearchBox },
  template: `<div>
      <NsComboSearchBox v-bind="$props" @search="testSearch('test')" @change="testChange('test')" @filter="testFilter('test')">
      <template slot="tooltip">
        <h6>Tooltip title</h6>
        <p>Tooltip description</p>
      </template>
    </NsComboSearchBox>
    <div class="mg-top-xlg">NsComboSearchBox component is ideal for selecting a value from a very large or undefined set of options. It enables the parent component to fetch options matching the user's search query and display them in the listbox.</div>
    <div class="mg-top-sm">NsComboSearchBox emits the following events:</div>
    <ul>
      <li><strong>search</strong> - emitted when the user types in the search box, with a debounce of 500ms. On this event the parent component should fetch the matching options</li>
      <li><strong>change</strong> - emitted when the user selects an option from the listbox</li>
      <li><strong>filter</strong> - emitted when the user types in or clears the search box</li>
    </ul>
  </div>`,
});

const fruitOptions = [
  "Apple",
  "Banana",
  "Cherry",
  "Date",
  "Elderberry",
  "Fig",
  "Grape",
  "Kiwi Fruit",
  "Lemon",
  "Lime",
  "Mango",
  "Orange",
  "Passion Fruit",
  "Raisin",
  "Satsuma",
  "Tangerine",
  "Ugli Fruit",
  "Watermelon",
].map((item) => {
  const nameVal = item.replace(/\W/, "_").toLowerCase();
  return {
    name: nameVal,
    label: item,
    value: nameVal,
  };
});

export const Default = Template.bind({});
Default.args = {
  resultsLimitReached: false,
  loadingResults: false,
  autoHighlight: false,
  disabled: false,
  invalidMessage: "",
  warnText: "",
  helperText: "",
  title: "Title",
  placeholder: "Search",
  highlight: "",
  value: "",
  light: true,
  clearFilterLabel: "Clear search",
  noResultsLabel: "No results",
  resultsLimitReachedLabel: "Keep typing to show more options",
  maxDisplayOptions: 100,
  tooltipAlignment: "start",
  tooltipDirection: "bottom",
  results: fruitOptions,
};
