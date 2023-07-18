module.exports = {
    branches: [
        { name: 'next' },
        { name: 'events-filter', prerelease: true },
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
