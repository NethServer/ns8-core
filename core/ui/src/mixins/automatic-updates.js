/*
 * Copyright (C) 2024 Nethesis S.r.l.
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

import to from "await-to-js";

/**
 * Shared helper to dispatch the `set-automatic-updates` cluster action.
 *
 * The component using this mixin must also mix in `UtilService` and
 * `TaskService` from `@nethserver/ns8-ui-lib`, since the method relies on
 * `createClusterTask`, `getUuid` and `getErrorMessage`.
 */
export default {
  methods: {
    /**
     * @param {Object} data - task payload, e.g. { apply_updates_is_active: false }
     *   or { instances: { "dokuwiki1": false } }.
     * @param {Object} options
     * @param {String} options.title - task title, shown in the completion toast.
     * @param {Function} [options.onCompleted] - called on task completion.
     * @param {Function} [options.onError] - called with an error message string.
     */
    async setAutomaticUpdates(data, { title, onCompleted, onError } = {}) {
      const taskAction = "set-automatic-updates";
      const eventId = this.getUuid();

      this.$root.$once(`${taskAction}-completed-${eventId}`, () => {
        if (onCompleted) {
          onCompleted();
        }
      });
      this.$root.$once(`${taskAction}-aborted-${eventId}`, (taskResult) => {
        console.error(`${taskAction} aborted`, taskResult);
        if (onError) {
          onError(this.$t("error.generic_error"));
        }
      });

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data,
          extra: {
            title: title || this.$t("action.set-automatic-updates"),
            eventId,
          },
        })
      );
      const err = res[0];
      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        if (onError) {
          onError(this.getErrorMessage(err));
        }
      }
    },
  },
};
