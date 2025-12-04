module.exports = {
  apps: [
    {
      name: "next-template",
      script: "npm",
      args: "start",
      exec_mode: "fork",
      cwd: "/home/nexttemplate/public_html",
      env: {
        PORT: "3200",
      },
    },
  ],
};
