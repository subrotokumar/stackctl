import { defineConfig } from 'vitepress'

// If the site is served from https://<user>.github.io/stackctl/ keep base as is.
// For a custom domain or user site, change it to '/'.
export default defineConfig({
  title: 'Stackctl',
  description: 'Create Java projects and run project tasks from your terminal.',
  base: '/stackctl/',
  lang: 'en-US',
  cleanUrls: true,
  lastUpdated: true,
  outDir: './dist',


  head: [['meta', { name: 'theme-color', content: '#e11d8d' }]],

  themeConfig: {
    nav: [
      { text: 'Guide', link: '/guide/what-is-stackctl', activeMatch: '/guide/' },
      { text: 'Reference', link: '/reference/cli', activeMatch: '/reference/' },
      { text: 'Examples', link: '/examples' },
    ],

    sidebar: [
      {
        text: 'Introduction',
        items: [
          { text: 'What is Stackctl?', link: '/guide/what-is-stackctl' },
          { text: 'Installation', link: '/guide/installation' },
          { text: 'Quick start', link: '/guide/quick-start' },
        ],
      },
      {
        text: 'Guide',
        items: [
          { text: 'Creating projects', link: '/guide/creating-projects' },
          { text: 'Task runner', link: '/guide/task-runner' },
          { text: 'Dependencies & aliases', link: '/guide/dependencies-aliases' },
          { text: 'Variables & templates', link: '/guide/variables-templates' },
          { text: 'Interactive selector', link: '/guide/interactive-selector' },
        ],
      },
      {
        text: 'Reference',
        items: [
          { text: 'CLI commands', link: '/reference/cli' },
          { text: 'Task file format', link: '/reference/taskfile' },
          { text: 'Troubleshooting', link: '/reference/troubleshooting' },
        ],
      },
      { text: 'Examples', link: '/examples' },
    ],

    socialLinks: [{ icon: 'github', link: 'https://github.com/subrotokumar/stackctl' }],

    search: { provider: 'local' },

    editLink: {
      pattern: 'https://github.com/subrotokumar/stackctl/edit/main/docs/:path',
      text: 'Edit this page on GitHub',
    },

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © Subroto Kumar',
    },
  },
})
