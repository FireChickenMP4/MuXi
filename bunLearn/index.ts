const server = Bun.serve({
  port: 6969,
  routes: {
    "/": () => new Response("Hi Mom"),
  },
});

console.log("Hello via Bun!");
