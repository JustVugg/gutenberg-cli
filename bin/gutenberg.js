#!/usr/bin/env node
import { main } from "../src/cli.js";

main(process.argv.slice(2)).catch((error) => {
  const message = error && error.message ? error.message : String(error);
  console.error(`gutenberg: ${message}`);
  process.exitCode = error && Number.isInteger(error.exitCode) ? error.exitCode : 1;
});
