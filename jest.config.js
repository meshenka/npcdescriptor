module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'jsdom',
  setupFilesAfterEnv: ['<rootDir>/frontend/src/setupTests.ts'],
  testMatch: ['**/frontend/src/**/*.test.tsx'],
  moduleNameMapper: {
    '\\.(png|jpg|webp|ttf|woff|woff2|svg)$': '<rootDir>/frontend/src/__mocks__/fileMock.js',
  },
};
