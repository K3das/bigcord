import { defineConfig } from 'vite';
import solidPlugin from 'vite-plugin-solid';

import { createHtmlPlugin } from 'vite-plugin-html'

export default defineConfig({
  plugins: [
    solidPlugin(),
    createHtmlPlugin({
      minify: {
        collapseWhitespace: true,
        keepClosingSlash: true,
        removeComments: false,
        removeRedundantAttributes: true,
        removeScriptTypeAttributes: true,
        removeStyleLinkTypeAttributes: true,
        useShortDoctype: true,
        minifyCSS: true,
      },
      template: 'index.html',
    }),
  ],
  build: {
    target: 'esnext',
    rollupOptions: {
      output: {
        compact: true,
        assetFileNames: "owo.[hash][extname]",
        chunkFileNames: "uwu.[hash].js",
        entryFileNames: "nya.[hash].js",
      },
    },
  },
  assetsInclude: [
    "public/*"
  ],
  base: ""
});
