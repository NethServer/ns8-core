//
// Copyright (C) 2026 Nethesis S.r.l.
// SPDX-License-Identifier: GPL-3.0-or-later
//

// Loki validates the pattern with RE2, which accepts patterns the backtracking
// JavaScript engine cannot run: "(.*)*x" costs 2^n and never returns. The regex
// API has no timeout and no interruption point, so matching runs here, where the
// caller can kill the worker past its deadline.

function highlightWorkerBody() {
  // leading timestamp and the tag added by the backend's LogQL line_format
  var PROCESS_TAG_PATTERN = /^(\S+\s+)(\[[^\]]*\])/;

  // past this many chunks a line is unreadable anyway
  var MAX_PARTS = 200;

  // flat [start, end, ...] offsets, so the reply carries no copy of the text
  function matchRanges(text, re) {
    var ranges = [];
    var parts = 0;
    var cursor = 0;
    var match;

    re.lastIndex = 0;
    while ((match = re.exec(text)) !== null) {
      // "()" matches empty at every position: without this the loop never ends
      if (!match[0]) {
        re.lastIndex++;
        if (re.lastIndex > text.length) {
          break;
        }
        continue;
      }
      // "(.)" matches every character: merge adjacent matches into one run
      if (match.index === cursor && ranges.length) {
        ranges[ranges.length - 1] = match.index + match[0].length;
      } else {
        if (match.index > cursor) {
          parts++;
        }
        ranges.push(match.index, match.index + match[0].length);
        parts++;
      }
      cursor = match.index + match[0].length;
      if (parts >= MAX_PARTS) {
        break;
      }
    }
    return ranges;
  }

  // matched per segment: a match may not span the timestamp/tag boundary
  function lineRanges(line, re) {
    var tag = PROCESS_TAG_PATTERN.exec(line);
    if (!tag) {
      return [[], [], matchRanges(line, re)];
    }
    return [
      matchRanges(tag[1], re),
      matchRanges(tag[2], re),
      matchRanges(line.slice(tag[0].length), re),
    ];
  }

  self.onmessage = function (event) {
    var request = event.data;
    var re;

    try {
      re = new RegExp(
        request.source,
        request.flags.indexOf("g") === -1 ? request.flags + "g" : request.flags
      );
    } catch (e) {
      // RE2 accepts syntax JavaScript rejects, e.g. "(?P<name>x)"
      self.postMessage({ id: request.id, error: "compile" });
      return;
    }

    var ranges = new Array(request.lines.length);
    var marks = 0;

    // from the end: the viewer scrolls to the bottom, so the budget covers what
    // is on screen
    for (var i = request.lines.length - 1; i >= 0; i--) {
      var line = lineRanges(request.lines[i], re);
      marks += (line[0].length + line[1].length + line[2].length) / 2;
      if (marks > request.maxMarks) {
        break;
      }
      ranges[i] = line;
    }
    self.postMessage({ id: request.id, ranges: ranges });
  };
}

// Blob because webpack 4 has no `new Worker(new URL(...))` and worker-loader is
// not a dependency. Null when unavailable: the caller then drops the highlighting
// rather than matching on the main thread.
export function createHighlightWorker() {
  try {
    const source = "(" + String(highlightWorkerBody) + ")()";
    const url = URL.createObjectURL(
      new Blob([source], { type: "application/javascript" })
    );
    const worker = new Worker(url);
    // the worker keeps its own reference to the script once constructed
    URL.revokeObjectURL(url);
    return worker;
  } catch (e) {
    return null;
  }
}
