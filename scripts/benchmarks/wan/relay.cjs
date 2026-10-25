#!/usr/bin/env node
// Transparent TCP relay with injected one-way propagation delay (per direction).
//
// Usage (legacy fixed RTT + optional jitter):
//   node relay.cjs <listenPort> <upstreamHost> <upstreamPort> <rttMs> [jitterMs] [seed]
//
// Usage (geographic band, jittered per connection via a seeded PRNG):
//   node relay.cjs <listenPort> <upstreamHost> <upstreamPort> --band <min:max> [seed]
//
// Model:
//   - Legacy mode: for every round trip the relay adds `rttMs` total (rttMs/2 on
//     each direction). `jitterMs` adds U(0, jitterMs) extra to each direction's
//     delay. `seed` makes a jittered run reproducible.
//   - Band mode (network component in full): each incoming TCP connection draws
//     ONE RTT from U(min, max) on a seeded PRNG, and both directions of that
//     connection use rtt/2 for every byte stream. Fresh connections draw fresh
//     values, so request/response round trips vary across the source's latency
//     class (same-region, cross-region, intercontinental). The seed keeps the
//     sequence of draws reproducible.
//
// Protocol-agnostic (HTTP/1.1, WebSocket, TLS, anything) because it only
// forwards raw TCP byte streams.
//
// Exits with code 1 if the upstream or listen socket fails at startup.

"use strict";

const net = require("net");

const args = process.argv.slice(2);
const LISTEN = args[0];
const HOST = args[1];
const PORT = args[2];
if (!LISTEN || !HOST || !PORT || !args[3]) {
  console.error(
    "usage: relay.cjs" +
      " <listenPort> <upstreamHost> <upstreamPort>" +
      " <rttMs> [jitterMs] [seed]  |  <listenPort> <upstreamHost> <upstreamPort>" +
      " --band <min:max> [seed]"
  );
  process.exit(2);
}

const listenPort = Number(LISTEN);
const upstreamHost = HOST;
const upstreamPort = Number(PORT);

let seed = 0;
const rand = () => {
  seed = (seed + 0x6d2b79f5) | 0;
  let t = Math.imul(seed ^ (seed >>> 15), 1 | seed);
  t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
  return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
};

// fixed one-way delay for legacy mode (band mode samples per connection)
let fixedOneWay;
let bandMin;
let bandMax;
let desc;

if (args[3] === "--band") {
  // jittered-band mode: one RTT draw per connection within the band
  const band = String(args[4] || "");
  const [min, max] = band.split("-").map(Number);
  if (!Number.isFinite(min) || !Number.isFinite(max) || max <= min) {
    console.error(`invalid --band "${band}" (expected "min-max", max > min)`);
    process.exit(2);
  }
  seed = Number(args[5] || 0) >>> 0;
  seed = seed || 20260910; // default deterministic seed, like the WAN runner
  bandMin = min;
  bandMax = max;
  desc = `band ${bandMin}-${bandMax}ms (RTT jittered per-connection) seed=${seed}`;
} else {
  const rttMs = Number(args[3]);
  const jitterMs = Number(args[4] || 0);
  seed = Number(args[5] || 0) >>> 0;
  fixedOneWay = rttMs / 2;
  desc =
    `rtt=${rttMs}ms (one-way=${fixedOneWay}ms)` +
    (jitterMs > 0 ? ` jitter=${jitterMs}ms` : "") +
    ` seed=${seed}`;
}

console.log(`relay listen ${listenPort} -> ${upstreamHost}:${upstreamPort} ${desc}`);

const server = net.createServer((local) => {
  // one RTT sample per connection (band mode); legacy mode uses the fixed value
  const oneWayMs = fixedOneWay == null
    ? (bandMin + rand() * (bandMax - bandMin)) / 2
    : fixedOneWay;
  console.log(`conn one-way=${oneWayMs.toFixed(1)}ms`);

  const wire = (a, b) => {
    const queued = [];
    let draining = false;

    const flush = () => {
      if (draining || queued.length === 0) return;
      draining = true;
      const msg = queued.shift();
      setTimeout(() => {
        try {
          if (!b.destroyed) b.write(msg.buf);
        } catch {
          /* ignore */
        }
        draining = false;
        flush();
      }, msg.delay);
      if (queued.length > 10000) {
        // runaway buffer (upstream not reading); drop the connection
        a.destroy();
        b.destroy();
      }
    };

    a.on("data", (chunk) => {
      queued.push({ buf: chunk, delay: oneWayMs });
      flush();
    });
    a.on("end", () => setTimeout(() => { try { b.end(); } catch { /* ignore */ } }, oneWayMs));
    a.on("error", () => { a.destroy(); b.destroy(); });
  };

  const upstream = net.connect({ host: upstreamHost, port: upstreamPort });
  upstream.on("connect", () => { wire(local, upstream); wire(upstream, local); });
  upstream.on("error", () => { local.destroy(); });
  local.on("error", () => { upstream.destroy(); });
});

server.on("error", (err) => {
  console.error(`relay error: ${err.message}`);
  process.exit(1);
});

server.listen(listenPort, "127.0.0.1", () => {
  console.log(`listening on 127.0.0.1:${listenPort}`);
});