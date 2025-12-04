module.exports = {
  apps: [
    {
      name: "dev-next-template",
      script: "npm",
      args: "start",
      exec_mode: "fork",
      cwd: "/home/nexttemplate/domains/dev.nexttemplate.com/public_html",
      env: {
        PORT: "3201",
      },
    },
  ],
};
