module.exports = {
  ci: {
    collect: {
      url: ['https://app:8888/supervision/workflow'],
      settings: {
        extraHeaders: JSON.stringify({Cookie: 'XSRF-TOKEN=abcde; Other=other'}),
        chromeFlags: "--disable-gpu --no-sandbox",
      },
    },
    assert: {
      assertions: {
        "categories:performance": ["warn", {"minScore": 0.9}],
        "categories:accessibility": ["error", {"minScore": 1}],
        "categories:best-practices": ["warn", {"minScore": 0.9}],
        "categories:seo": ["warn", {"minScore": 0.9}],
      },
    },
    upload: {
      target: 'temporary-public-storage',
    },
  },
};
