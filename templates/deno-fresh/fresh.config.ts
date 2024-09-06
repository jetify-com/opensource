import { defineConfig } from "$fresh/server.ts";

// Use port 8080 to match the default port used by Jetify Cloud
export default defineConfig({ server: { port: 8080 } });
