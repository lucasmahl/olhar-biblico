import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";
import fs from "fs";
import { config } from "./src/lib/config";

const {
  PRODUCTION = false,
  CERT_FILE = "",
  KEY_FILE = "",
  PORT = 3000,
} = config() ?? {};

let serverConfiguration = {};
if (PRODUCTION === "TRUE") {
  serverConfiguration = {
    host: "olharbiblico.com",
    port: PORT,
    https: {
      key: fs.readFileSync(KEY_FILE),
      cert: fs.readFileSync(CERT_FILE),
    },
  };
} else {
  serverConfiguration = {
    host: "localhost",
    port: PORT,
  };
}

export default defineConfig({
  plugins: [sveltekit()],
  server: serverConfiguration,
});
