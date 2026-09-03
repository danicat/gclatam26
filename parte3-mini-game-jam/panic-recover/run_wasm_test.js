"use strict";

const fs = require("fs");
const path = require("path");

if (process.argv.length < 3) {
  console.error("usage: node run_wasm_test.js [wasm binary] [arguments]");
  process.exit(1);
}

globalThis.require = require;
globalThis.fs = fs;
globalThis.path = path;
globalThis.TextEncoder = require("util").TextEncoder;
globalThis.TextDecoder = require("util").TextDecoder;
globalThis.performance ??= { now: () => Date.now() };
globalThis.crypto ??= require("crypto");

// Mock browser objects for Ebitengine WASM headless testing
const dummyElement = {
  style: {},
  addEventListener: () => {},
  removeEventListener: () => {},
  appendChild: () => {},
  setAttribute: () => {},
  removeAttribute: () => {},
  setPointerCapture: () => {},
  releasePointerCapture: () => {},
  getContext: () => ({
    getExtension: () => null,
    getParameter: () => 0,
    enable: () => {},
    disable: () => {},
  }),
  getBoundingClientRect: () => ({ left: 0, top: 0, width: 640, height: 360 }),
};

globalThis.screen = { width: 1920, height: 1080 };
globalThis.requestAnimationFrame = (cb) => setTimeout(cb, 16);

globalThis.Document = function() {};
globalThis.Document.prototype = {};
Object.defineProperty(globalThis.Document.prototype, "hidden", {
  get: () => false,
  configurable: true,
});

globalThis.document = {
  createElement: () => dummyElement,
  head: dummyElement,
  body: dummyElement,
  documentElement: dummyElement,
  addEventListener: () => {},
  removeEventListener: () => {},
  getElementById: () => dummyElement,
  hasFocus: () => true,
};
globalThis.window = globalThis;
globalThis.addEventListener = () => {};
globalThis.removeEventListener = () => {};
globalThis.matchMedia = () => ({ matches: false, addEventListener: () => {}, removeEventListener: () => {} });
// navigator already exists in Node 22

require("./web/wasm_exec.js");

const go = new Go();
go.argv = process.argv.slice(2);
// Prune env to avoid env size overflow
go.env = { TMPDIR: require("os").tmpdir(), PATH: process.env.PATH || "" };
go.exit = process.exit;

WebAssembly.instantiate(fs.readFileSync(process.argv[2]), go.importObject).then((result) => {
  process.on("exit", (code) => {
    if (code === 0 && !go.exited) {
      go._pendingEvent = { id: 0 };
      go._resume();
    }
  });
  return go.run(result.instance);
}).catch((err) => {
  console.error(err);
  process.exit(1);
});
