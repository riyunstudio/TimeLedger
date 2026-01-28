import { c as defineEventHandler, u as useRuntimeConfig, g as getRouterParam, r as readBody, e as getRequestHeaders } from '../../_/nitro.mjs';
import 'node:http';
import 'node:https';
import 'node:events';
import 'node:buffer';
import 'node:fs';
import 'node:path';
import 'node:crypto';
import 'node:url';

const ____path_ = defineEventHandler(async (event) => {
  useRuntimeConfig();
  const apiBase = process.env.API_BASE_URL || "http://localhost:8888/api/v1";
  const path = getRouterParam(event, "path");
  const url = `${apiBase}/${path}`;
  try {
    const response = await fetch(url, {
      method: event.method,
      headers: {
        ...getRequestHeaders(event),
        "Content-Type": "application/json",
        "Accept": "application/json"
      },
      body: event.method !== "GET" && event.method !== "HEAD" ? await readBody(event) : void 0
    });
    const data = await response.json();
    return data;
  } catch (error) {
    return {
      code: 500,
      message: "API request failed",
      error: error instanceof Error ? error.message : "Unknown error"
    };
  }
});

export { ____path_ as default };
//# sourceMappingURL=_...path_.mjs.map
