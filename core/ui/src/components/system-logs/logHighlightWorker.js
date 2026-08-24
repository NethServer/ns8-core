//
// Copyright (C) 2026 Nethesis S.r.l.
// SPDX-License-Identifier: GPL-3.0-or-later
//

// The search pattern comes from the user and runs through the JavaScript regex
// engine, which backtracks. Loki validates with RE2, which simulates the
// automaton instead and therefore accepts patterns no backtracking engine can
// run: "(.*)*x" costs 2^n and never returns on a real log line. There is no
// timeout in the regex API and no interruption point, so running it on the main
// thread freezes the tab beyond recovery.
//
// So it runs here instead, where the caller can enforce a deadline by killing
// the worker. Terminating is the only reliable way to stop a regex mid-match.
//
// The matching logic below must stay identical to what LogHighlightMark used to
// do inline, including running the pattern separately on each segment of a line:
// a match may not span the timestamp/tag boundary, and matching the whole line
// then clipping would highlight different text.

function highlightWorkerBody() {
  // leading timestamp followed by the "[node_id:module_id:syslog_id]" tag
  // added by the backend's LogQL line_format
  var PROCESS_TAG_PATTERN = /^(\S+\s+)(\[[^\]]*\])/;

  // past this many alternating match/non-match chunks a line stops being
  // readable anyway, so splitting further only costs render time
  var MAX_PARTS = 200;

  // Returns match offsets as a flat [start, end, start, end, ...] list, so the
  // reply carries numbers instead of a second copy of the whole log text.
  function matchRanges(text, re) {
    var ranges = [];
    var parts = 0;
    var cursor = 0;
    var match;

    re.lastIndex = 0;
    while ((match = re.exec(text)) !== null) {
      // a pattern like "()" matches nothing at every position: skipping it
      // avoids both an endless loop and one part per character
      if (!match[0]) {
        re.lastIndex++;
        if (re.lastIndex > text.length) {
          break;
        }
        continue;
      }
      // parts are pushed as [optional gap][match], so whatever sits at the end
      // is always a match: touching it means this one continues the same run
      if (match.index === cursor && ranges.length) {
        // a pattern like "(.)" matches every character: merging adjacent
        // matches keeps one mark per run instead of one per character
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
        // pattern matching a huge number of scattered fragments: leave the
        // tail unhighlighted rather than flood the renderer
        break;
      }
    }
    return ranges;
  }

  // one entry per line: ranges for the timestamp, the tag, and the message
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

    var ranges = [];
    var marks = 0;

    for (var i = 0; i < request.lines.length; i++) {
      var line = lineRanges(request.lines[i], re);
      marks += (line[0].length + line[1].length + line[2].length) / 2;
      if (marks > request.maxMarks) {
        self.postMessage({ id: request.id, error: "too_many" });
        return;
      }
      ranges.push(line);
    }
    self.postMessage({ id: request.id, ranges: ranges });
  };
}

// Built from a Blob because the toolchain is webpack 4, which has no support for
// `new Worker(new URL(...))`, and worker-loader is not a dependency. Returns
// null when workers or blob URLs are unavailable, so the caller can fall back to
// showing the lines unhighlighted -- never to matching on the main thread.
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
