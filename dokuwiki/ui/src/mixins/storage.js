//// FILE COPIED FROM CORE

export default {
  name: "StorageService",
  methods: {
    getFromStorage(prop) {
      return JSON.parse(localStorage.getItem(prop));
    },
    saveToStorage(prop, object) {
      localStorage.setItem(prop, JSON.stringify(object));
    },
    deleteFromStorage(prop) {
      localStorage.removeItem(prop);
    },
  },
};
