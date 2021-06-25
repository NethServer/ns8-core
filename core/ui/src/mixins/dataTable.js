import UtilService from "@/mixins/util";

export default {
  name: "DataTableService",
  mixins: [UtilService],
  data() {
    return {
      pageStart: 0,
      pageNumber: 0,
      pageLength: 0,
    };
  },
  computed: {
    pagination() {
      return {
        numberOfItems: this.tableRows.length,
        pageSizes: [10, 25, 50, 100],
      };
    },
    tablePage() {
      if (this.pagination) {
        return this.tableRows.slice(
          this.pageStart,
          this.pageStart + this.pageLength
        );
      } else {
        return this.tableRows;
      }
    },
  },
  methods: {
    paginateTable(ev) {
      this.pageStart = ev.start - 1;
      this.pageNumber = ev.page;
      this.pageLength = ev.length;
    },
    sortTable(ev) {
      const order = ev.order;

      if (order === "none") {
        return;
      }

      const propertyToSort = this.tableColumns[ev.index];
      this.tableRows.sort(this.sortByProperty(propertyToSort));

      if (ev.order === "descending") {
        this.tableRows.reverse();
      }
    },
  },
};
