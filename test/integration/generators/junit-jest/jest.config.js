module.exports = {
  testEnvironment: "node",
  reporters: [["jest-junit", { outputDirectory: process.env.JEST_JUNIT_OUTPUT_DIR, outputName: "junit.xml" }]],
};
