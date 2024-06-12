module.exports = {
    content: ["./public/bundle.js"],
    theme: {
        extend: {
            colors: {
                "primary-green": "var(--primary-green)",
                "light-green": "var(--light-green)",
                "house-green": "var(--house-green)",
                "tea-green": "var(--tea-green)",
                "neutral-warm": "var(--neutral-warm)",
            },
        },
    },
    variants: {
        extend: {},
    },
    plugins: [],
};
