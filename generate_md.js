const fs = require('fs');
const endpoints = JSON.parse(fs.readFileSync('endpoints_utf8.json', 'utf8'));

let markdown = '# API 端點完整對照表\n\n**最後更新**：2026年1月31日\n**來源**：swagger.json\n\n---\n\n';

const categories = {};

endpoints.forEach(ep => {
    const tag = ep.tags[0] || 'Uncategorized';
    if (!categories[tag]) categories[tag] = [];
    categories[tag].push(ep);
});

for (const category in categories) {
    markdown += `## ${category}\n\n`;
    markdown += '| Method | Endpoint | Summary |\n';
    markdown += '|:---|:---|:---|\n';
    categories[category].forEach(ep => {
        markdown += `| ${ep.method} | \`${ep.path}\` | ${ep.summary} |\n`;
    });
    markdown += '\n';
}

fs.writeFileSync('API_REFERENCE_NEW.md', markdown);
