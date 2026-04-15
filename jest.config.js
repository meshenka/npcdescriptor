module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'jsdom',
  setupFilesAfterEnv: ['<rootDir>/frontend/src/setupTests.ts'],
  testMatch: ['**/frontend/src/**/*.test.tsx'],
};
