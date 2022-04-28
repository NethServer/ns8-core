import { NsDataTable, NsEmptyState } from "@nethserver/ns8-ui-lib";
import Edit20 from "@carbon/icons-vue/es/edit/20";
import TrashCan20 from "@carbon/icons-vue/es/trash-can/20";

export default {
  title: "Components/NsDataTable",
  component: NsDataTable,
  argTypes: {
    rowSize: {
      options: ["compact", "short", "standard", "tall", ""],
      control: { type: "radio" },
    },
    editUser: { action: "editUserAction" },
    deleteUser: { action: "deleteUserAction" },
  },
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsDataTable, NsEmptyState },
  data() {
    return {
      tablePage: [],
      Edit20,
      TrashCan20,
    };
  },
  template:
    '<cv-tile light>\
      <NsDataTable v-bind="$props" @updatePage="tablePage = $event">\
        <template slot="empty-state">\
          <NsEmptyState title="No data">\
          </NsEmptyState>\
        </template>\
        <template slot="data">\
          <cv-data-table-row\
            v-for="(row, rowIndex) in tablePage"\
            :key="`${rowIndex}`"\
            :value="`${rowIndex}`"\
          >\
            <cv-data-table-cell>{{ row.fullName }}</cv-data-table-cell>\
            <cv-data-table-cell>{{ row.email }}</cv-data-table-cell>\
            <cv-data-table-cell>{{ row.phone }}</cv-data-table-cell>\
            <cv-data-table-cell>{{ row.country }}</cv-data-table-cell>\
            <cv-data-table-cell\
              class="table-overflow-menu-cell"\
            >\
              <cv-overflow-menu flip-menu class="table-overflow-menu">\
                <cv-overflow-menu-item @click="editUser(row)">\
                  <NsMenuItem :icon="Edit20" label="Edit" />\
                </cv-overflow-menu-item>\
                <cv-overflow-menu-item danger @click="deleteUser(row)">\
                  <NsMenuItem :icon="TrashCan20" label="Delete" />\
                </cv-overflow-menu-item>\
              </cv-overflow-menu>\
            </cv-data-table-cell>\
          </cv-data-table-row>\
        </template>\
      </NsDataTable>\
      <div style="margin-top: 2rem">\
        See <a href="https://vue.carbondesignsystem.com/?path=/story/components-cvdatatable--default" target="_blank">https://vue.carbondesignsystem.com/?path=/story/components-cvdatatable--default</a> for other examples\
      </div>\
    </cv-tile>',
});

const users = [
  {
    fullName: "Anton Peltola",
    email: "anton.peltola@example.com",
    phone: "06-450-409",
    country: "Finland",
  },
  {
    fullName: "Mohini Overeem",
    email: "mohini.overeem@example.com",
    phone: "(395)-518-0029",
    country: "Netherlands",
  },
  {
    fullName: "Mason Lewis",
    email: "mason.lewis@example.com",
    phone: "(936)-975-9578",
    country: "New Zealand",
  },
  {
    fullName: "Loane Clement",
    email: "loane.clement@example.com",
    phone: "05-80-51-35-04",
    country: "France",
  },
  {
    fullName: "Linnea Rintala",
    email: "linnea.rintala@example.com",
    phone: "03-339-837",
    country: "Finland",
  },
  {
    fullName: "Julie Jean",
    email: "julie.jean@example.com",
    phone: "05-27-21-47-50",
    country: "France",
  },
  {
    fullName: "Jessica Jackson",
    email: "jessica.jackson@example.com",
    phone: "(838)-538-3251",
    country: "United States",
  },
  {
    fullName: "Ryder White",
    email: "ryder.white@example.com",
    phone: "(183)-200-9152",
    country: "New Zealand",
  },
  {
    fullName: "Connor Robinson",
    email: "connor.robinson@example.com",
    phone: "(965)-217-4178",
    country: "Switzerland",
  },
  {
    fullName: "Jim Bradley",
    email: "jim.bradley@example.com",
    phone: "00-3880-6902",
    country: "Australia",
  },
  {
    fullName: "Arnold Wells",
    email: "arnold.wells@example.com",
    phone: "031-952-6372",
    country: "Ireland",
  },
  {
    fullName: "Adeline Gaillard",
    email: "adeline.gaillard@example.com",
    phone: "075 106 31 13",
    country: "Switzerland",
  },
  {
    fullName: "Theo Hunter",
    email: "theo.hunter@example.com",
    phone: "071-302-7111",
    country: "Ireland",
  },
  {
    fullName: "Beatrice Spencer",
    email: "beatrice.spencer@example.com",
    phone: "00-4814-9335",
    country: "Australia",
  },
  {
    fullName: "Suzy Warren",
    email: "suzy.warren@example.com",
    phone: "011-103-1348",
    country: "United Kingdom",
  },
];

const sharedProps = {
  autoWidth: false,
  borderless: false,
  overflowMenu: true,
  rowSize: "standard",
  searchLabel: "Search",
  searchPlaceholder: "Search",
  searchClearLabel: "Clear search",
  sortable: true,
  title: "",
  columns: ["Full name", "Email", "Phone", "Country"],
  rawColumns: ["fullName", "email", "phone", "country"],
  zebra: false,
  stickyHeader: false,
  helperText: "",
  staticWidth: false,
  isSearchable: true,
  skeletonRows: 5,
  noSearchResultsLabel: "No results",
  noSearchResultsDescription: "Try changing your search query",
  errorTitle: "Cannot retrieve table data",
  errorDescription: "Something went wrong",
  itemsPerPageLabel: "Items per page:",
  rangeOfTotalItemsLabel: "{range} of {total} items",
  ofTotalPagesLabel: "of {total} pages",
  backwardText: "Prev page",
  forwardText: "Next page",
  pageNumberLabel: "Page number:",
};

export const Default = Template.bind({});
Default.args = {
  pageSizes: [5, 10, 25, 50, 100],
  isLoading: false,
  isErrorShown: false,
  allRows: users,
  ...sharedProps,
};

export const EmptyState = Template.bind({});
EmptyState.args = {
  pageSizes: [5, 10, 25, 50, 100],
  isLoading: false,
  isErrorShown: false,
  allRows: [],
  ...sharedProps,
};

export const NoPagination = Template.bind({});
NoPagination.args = {
  pageSizes: [], // no pagination
  isLoading: false,
  isErrorShown: false,
  allRows: users,
  ...sharedProps,
};

export const Loading = Template.bind({});
Loading.args = {
  pageSizes: [5, 10, 25, 50, 100],
  isLoading: true,
  isErrorShown: false,
  allRows: [],
  ...sharedProps,
};

export const Error = Template.bind({});
Error.args = {
  pageSizes: [5, 10, 25, 50, 100],
  isLoading: false,
  isErrorShown: true,
  allRows: [],
  ...sharedProps,
};
