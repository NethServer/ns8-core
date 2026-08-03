<!--
  Copyright (C) 2026 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div ref="container" class="terminal-container"></div>
</template>

<script>
// The 5.3 series rather than the renamed @xterm/xterm 5.5: webpack 4, which
// vue-cli 4 pins, cannot parse the syntax of the newer bundle.
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import "xterm/css/xterm.css";

export default {
  name: "XtermPane",
  props: {
    // Open websocket. Binary frames carry terminal bytes, text frames control
    // messages. AttachAddon is deliberately not used: it is bidirectional by
    // default, so pairing it with our own onData handler sends every keystroke
    // twice.
    //
    // Becomes null once the session ends: the pane outlives it so the operator
    // keeps the transcript and the closing notice on screen.
    socket: {
      type: WebSocket,
      default: null,
    },
  },
  data() {
    return {
      term: null,
      fitAddon: null,
      observer: null,
      lastSize: { rows: 0, cols: 0 },
    };
  },
  mounted() {
    this.term = new Terminal({
      fontFamily: "'IBM Plex Mono', 'Menlo', monospace",
      fontSize: 13,
      cursorBlink: true,
      scrollback: 5000,
      theme: this.terminalTheme(),
    });

    this.fitAddon = new FitAddon();
    this.term.loadAddon(this.fitAddon);
    this.term.open(this.$refs.container);

    this.term.onData(this.onTerminalData);
    this.term.onResize(this.onTerminalResize);

    this.observer = new ResizeObserver(this.fit);
    this.observer.observe(this.$refs.container);

    if (this.socket) {
      this.socket.binaryType = "arraybuffer";
      this.socket.addEventListener("message", this.onSocketMessage);
    }

    this.fit();
    this.term.focus();
  },
  beforeDestroy() {
    if (this.observer) {
      this.observer.disconnect();
    }
    if (this.socket) {
      this.socket.removeEventListener("message", this.onSocketMessage);
    }
    if (this.term) {
      this.term.dispose();
    }
  },
  methods: {
    terminalTheme() {
      // Carbon gray 100 so the pane matches the dark theme of the shell.
      return {
        background: "#161616",
        foreground: "#f4f4f4",
        cursor: "#f4f4f4",
        selectionBackground: "#525252",
      };
    },
    fit() {
      if (!this.fitAddon) {
        return;
      }
      try {
        this.fitAddon.fit();
      } catch (error) {
        // The container can be measured as zero while the view transitions.
        return;
      }
    },
    onTerminalData(data) {
      if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
        return;
      }
      // Encode explicitly: the relay expects bytes, and a JS string would lose
      // the distinction between control frames and terminal input.
      this.socket.send(new TextEncoder().encode(data));
    },
    onTerminalResize({ rows, cols }) {
      if (
        !this.socket ||
        this.socket.readyState !== WebSocket.OPEN ||
        (rows === this.lastSize.rows && cols === this.lastSize.cols)
      ) {
        return;
      }
      this.lastSize = { rows, cols };
      this.socket.send(JSON.stringify({ type: "resize", rows, cols }));
    },
    onSocketMessage(event) {
      if (typeof event.data === "string") {
        // Control frames are handled by the parent view.
        return;
      }
      // Write the bytes as they arrive: decoding to a string first would break
      // a multi-byte character split across two frames.
      this.term.write(new Uint8Array(event.data));
    },
    // Used by the parent to replay bytes that arrived before this pane was
    // mounted: the login prompt is the very first thing the server writes.
    writeBytes(data) {
      if (this.term) {
        this.term.write(new Uint8Array(data));
      }
    },
    writeNotice(message) {
      if (this.term) {
        this.term.write(`\r\n\x1b[31m${message}\x1b[0m\r\n`);
      }
    },
    currentSize() {
      return {
        rows: this.term ? this.term.rows : 24,
        cols: this.term ? this.term.cols : 80,
      };
    },
  },
};
</script>

<style scoped lang="scss">
.terminal-container {
  height: 60vh;
  min-height: 20rem;
  padding: 0.5rem;
  background-color: #161616;
}
</style>
