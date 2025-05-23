import { createMDX } from 'fumadocs-mdx/next';

const withMDX = createMDX();

/** @type {import('next').NextConfig} */
const config = {
  reactStrictMode: true,
  async headers() {
    return [
      {
        source: '/mowalang.wasm',
        headers: [
          {
            key: 'Content-Type',
            value: 'application/wasm',
          },
        ],
      },
      {
        source: '/wasm_exec.js',
        headers: [
          {
            key: 'Content-Type',
            value: 'text/javascript',
          },
        ],
      },
    ];
  },
};

export default withMDX(config);
